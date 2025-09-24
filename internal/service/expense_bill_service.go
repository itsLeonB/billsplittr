package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/mapper"
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
)

type expenseBillServiceImpl struct {
	transactor crud.Transactor
	billRepo   repository.ExpenseBillRepository
	logger     ezutil.Logger
	taskQueue  repository.TaskQueue
}

func NewExpenseBillService(
	transactor crud.Transactor,
	billRepo repository.ExpenseBillRepository,
	logger ezutil.Logger,
	taskQueue repository.TaskQueue,
) ExpenseBillService {
	return &expenseBillServiceImpl{
		transactor,
		billRepo,
		logger,
		taskQueue,
	}
}

func (ebs *expenseBillServiceImpl) Save(ctx context.Context, req dto.NewExpenseBillRequest) (dto.ExpenseBillResponse, error) {
	newBill := entity.ExpenseBill{
		PayerProfileID:   req.PayerProfileID,
		CreatorProfileID: req.CreatorProfileID,
		ImageName:        req.Filename,
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
	if err := ebs.taskQueue.Ping(); err != nil {
		return err
	}

	spec := crud.Specification[entity.ExpenseBill]{}
	spec.DeletedFilter = crud.ExcludeDeleted
	bills, err := ebs.billRepo.FindAll(ctx, spec)
	if err != nil {
		return err
	}

	validObjectKeys := ezutil.MapSlice(bills, func(eb entity.ExpenseBill) string { return eb.ImageName })
	task, err := entity.NewTask(entity.OrphanedBillCleanupTask{BillObjectKeys: validObjectKeys})
	if err != nil {
		return err
	}

	return ebs.taskQueue.Enqueue(ctx, task)
}

func (s *expenseBillServiceImpl) getBySpec(ctx context.Context, spec crud.Specification[entity.ExpenseBill]) (entity.ExpenseBill, error) {
	bill, err := s.billRepo.FindFirst(ctx, spec)
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
