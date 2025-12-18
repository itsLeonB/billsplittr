package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/client"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/mapper"
	"github.com/itsLeonB/billsplittr/internal/message"
	"github.com/itsLeonB/billsplittr/internal/pkg/config"
	"github.com/itsLeonB/billsplittr/internal/pkg/logger"
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/meq"
	"github.com/itsLeonB/ungerr"
)

type expenseBillServiceImpl struct {
	transactor   crud.Transactor
	billRepo     repository.ExpenseBillRepository
	cleanupQueue meq.TaskQueue[message.OrphanedBillCleanup]
	bucketName   string
	ocr          client.OCRClient
}

func NewExpenseBillService(
	transactor crud.Transactor,
	billRepo repository.ExpenseBillRepository,
	cleanupQueue meq.TaskQueue[message.OrphanedBillCleanup],
	bucketName string,
	ocr client.OCRClient,
) ExpenseBillService {
	return &expenseBillServiceImpl{
		transactor,
		billRepo,
		cleanupQueue,
		bucketName,
		ocr,
	}
}

func (ebs *expenseBillServiceImpl) Save(ctx context.Context, req dto.NewExpenseBillRequest) (dto.ExpenseBillResponse, error) {
	expenseID := uuid.NullUUID{}
	if req.GroupExpenseID != uuid.Nil {
		expenseID = uuid.NullUUID{
			UUID:  req.GroupExpenseID,
			Valid: true,
		}
	}

	newBill := entity.ExpenseBill{
		PayerProfileID:   req.PayerProfileID,
		CreatorProfileID: req.CreatorProfileID,
		GroupExpenseID:   expenseID,
		ImageName:        req.Filename,
		Status:           appconstant.PendingBill,
	}

	insertedBill, err := ebs.billRepo.Insert(ctx, newBill)
	if err != nil {
		return dto.ExpenseBillResponse{}, err
	}

	return mapper.ExpenseBillToResponse(insertedBill), nil
}

func (ebs *expenseBillServiceImpl) GetAllCreated(ctx context.Context, creatorProfileID uuid.UUID) ([]dto.ExpenseBillResponse, error) {
	spec := crud.Specification[entity.ExpenseBill]{}
	spec.Model.CreatorProfileID = creatorProfileID
	spec.DeletedFilter = crud.ExcludeDeleted

	bills, err := ebs.billRepo.FindAll(ctx, spec)
	if err != nil {
		return nil, err
	}

	return ezutil.MapSlice(bills, mapper.ExpenseBillToResponse), nil
}

func (ebs *expenseBillServiceImpl) Get(ctx context.Context, id uuid.UUID) (dto.ExpenseBillResponse, error) {
	spec := crud.Specification[entity.ExpenseBill]{}
	spec.Model.ID = id
	spec.DeletedFilter = crud.ExcludeDeleted

	bill, err := ebs.getBySpec(ctx, spec)
	if err != nil {
		return dto.ExpenseBillResponse{}, err
	}

	return mapper.ExpenseBillToResponse(bill), nil
}

func (ebs *expenseBillServiceImpl) Delete(ctx context.Context, id, profileID uuid.UUID) error {
	return ebs.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := crud.Specification[entity.ExpenseBill]{}
		spec.Model.ID = id
		spec.Model.CreatorProfileID = profileID
		spec.ForUpdate = true
		spec.DeletedFilter = crud.ExcludeDeleted

		bill, err := ebs.getBySpec(ctx, spec)
		if err != nil {
			return err
		}

		bill.DeletedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

		_, err = ebs.billRepo.Update(ctx, bill)

		return err
	})
}

func (ebs *expenseBillServiceImpl) EnqueueCleanup(ctx context.Context) error {
	spec := crud.Specification[entity.ExpenseBill]{}
	bills, err := ebs.billRepo.FindAll(ctx, spec)
	if err != nil {
		return err
	}

	validObjectKeys := ezutil.MapSlice(bills, func(eb entity.ExpenseBill) string { return eb.ImageName })

	logger.Global.Infof("obtained object keys from DB:\n%s", strings.Join(validObjectKeys, "\n"))

	task := message.OrphanedBillCleanup{
		BillObjectKeys: validObjectKeys,
		BucketName:     ebs.bucketName,
	}

	return ebs.cleanupQueue.Enqueue(ctx, config.AppName, task)
}

func (ebs *expenseBillServiceImpl) ExtractBillText(ctx context.Context, msg message.ExpenseBillUploaded) (string, error) {
	var extractedText string
	err := ebs.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := crud.Specification[entity.ExpenseBill]{}
		spec.Model.ID = msg.ID
		spec.ForUpdate = true
		bill, err := ebs.getBySpec(ctx, spec)
		if err != nil {
			return err

		}

		text, err := ebs.ocr.ExtractFromURI(ctx, msg.URI)
		if err != nil {
			bill.Status = appconstant.FailedExtracting
			_, statusErr := ebs.billRepo.Update(ctx, bill)
			if statusErr != nil {
				return errors.Join(err, statusErr)
			}
			return err
		}

		extractedText = text
		bill.ExtractedText = text
		bill.Status = appconstant.ExtractedBill
		_, err = ebs.billRepo.Update(ctx, bill)
		return err
	})
	return extractedText, err
}

func (ebs *expenseBillServiceImpl) getBySpec(ctx context.Context, spec crud.Specification[entity.ExpenseBill]) (entity.ExpenseBill, error) {
	bill, err := ebs.billRepo.FindFirst(ctx, spec)
	if err != nil {
		return entity.ExpenseBill{}, err
	}
	if bill.IsZero() {
		return entity.ExpenseBill{}, ungerr.NotFoundError(fmt.Sprintf("expense bill with ID %s is not found", spec.Model.ID))
	}
	if bill.IsDeleted() {
		return entity.ExpenseBill{}, ungerr.UnprocessableEntityError(fmt.Sprintf("expense bill with ID %s is deleted", spec.Model.ID))
	}
	return bill, nil
}
