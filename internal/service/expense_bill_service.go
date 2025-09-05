package service

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
	"github.com/rotisserie/eris"
)

type expenseBillServiceImpl struct {
	billRepo    repository.ExpenseBillRepository
	storageRepo repository.StorageRepository
	bucketName  string
	logger      ezutil.Logger
}

func NewExpenseBillService(
	billRepo repository.ExpenseBillRepository,
	storageRepo repository.StorageRepository,
	bucketName string,
	logger ezutil.Logger,
) ExpenseBillService {
	if bucketName == "" {
		panic("bucket name cannot be empty")
	}
	return &expenseBillServiceImpl{
		billRepo,
		storageRepo,
		bucketName,
		logger,
	}
}

func (s *expenseBillServiceImpl) Upload(ctx context.Context, req *dto.UploadBillRequest) (uuid.UUID, error) {
	// Validate request
	if err := s.validateUploadRequest(req); err != nil {
		return uuid.Nil, err
	}

	// Generate object key
	objectKey := s.generateObjectKey(req.CreatorProfileID, req.Filename)

	// Upload to storage
	storageReq := &entity.StorageUploadRequest{
		Data:        req.ImageData,
		ContentType: req.ContentType,
		Filename:    req.Filename,
		BucketName:  s.bucketName,
		ObjectKey:   objectKey,
	}

	if _, err := s.storageRepo.Upload(ctx, storageReq); err != nil {
		return uuid.Nil, eris.Wrap(err, appconstant.ErrStorageUploadFailed)
	}

	// Create bill entity
	bill := entity.ExpenseBill{
		PayerProfileID:   req.PayerProfileID,
		CreatorProfileID: req.CreatorProfileID,
		ImageName:        objectKey,
	}

	// Save to database
	savedBill, err := s.billRepo.Insert(ctx, bill)
	if err != nil {
		// Try to clean up uploaded file
		_ = s.storageRepo.Delete(ctx, s.bucketName, objectKey)
		return uuid.Nil, eris.Wrap(err, "failed to save bill to database")
	}

	return savedBill.ID, nil
}

func (s *expenseBillServiceImpl) Get(ctx context.Context, billID uuid.UUID, profileID uuid.UUID) (dto.ExpenseBillResponse, error) {
	bill, err := s.getByID(ctx, billID, profileID)
	if err != nil {
		return dto.ExpenseBillResponse{}, err
	}

	return dto.ExpenseBillResponse{
		ID:               bill.ID,
		CreatorProfileID: bill.CreatorProfileID,
		PayerProfileID:   bill.PayerProfileID,
		GroupExpenseID:   bill.GroupExpenseID.UUID,
		Filename:         bill.ImageName,
		CreatedAt:        bill.CreatedAt,
		UpdatedAt:        bill.UpdatedAt,
		DeletedAt:        bill.DeletedAt.Time,
	}, nil
}

func (s *expenseBillServiceImpl) Delete(ctx context.Context, billID uuid.UUID, profileID uuid.UUID) error {
	bill, err := s.getByID(ctx, billID, profileID)
	if err != nil {
		return err
	}

	// Delete from storage
	if err := s.storageRepo.Delete(ctx, s.bucketName, bill.ImageName); err != nil {
		// Log error but don't fail the operation
		// You might want to use a proper logger here
		s.logger.Warnf("failed to delete file from storage: %v\n", err)
	}

	// Delete from database
	spec := entity.ExpenseBill{}
	spec.ID = billID
	return s.billRepo.Delete(ctx, spec)
}

func (s *expenseBillServiceImpl) getByID(ctx context.Context, id, profileID uuid.UUID) (entity.ExpenseBill, error) {
	spec := crud.Specification[entity.ExpenseBill]{}
	spec.Model.ID = id
	spec.Model.CreatorProfileID = profileID
	bill, err := s.billRepo.FindFirst(ctx, spec)
	if err != nil {
		return entity.ExpenseBill{}, err
	}
	if bill.IsZero() {
		return entity.ExpenseBill{}, ungerr.NotFoundError(fmt.Sprintf("expense bill with ID %s is not found", id))
	}
	if bill.IsDeleted() {
		return entity.ExpenseBill{}, ungerr.UnprocessableEntityError(fmt.Sprintf("expense bill with ID %s is deleted", id))
	}
	return bill, nil
}

func (s *expenseBillServiceImpl) validateUploadRequest(req *dto.UploadBillRequest) error {
	if req.CreatorProfileID == uuid.Nil {
		return ungerr.BadRequestError("creator profile ID is required")
	}

	if len(req.ImageData) == 0 {
		return ungerr.BadRequestError("image data is required")
	}

	if len(req.ImageData) > appconstant.MaxFileSize {
		return ungerr.BadRequestError(appconstant.ErrFileTooLarge)
	}

	if _, ok := appconstant.AllowedContentTypes[req.ContentType]; !ok {
		return ungerr.BadRequestError(appconstant.ErrInvalidFileType)
	}

	return nil
}

func (s *expenseBillServiceImpl) generateObjectKey(creatorID uuid.UUID, filename string) string {
	ext := filepath.Ext(filename)
	timestamp := time.Now().Format("2006/01/02")
	return fmt.Sprintf("bills/%s/%s/%s%s", timestamp, creatorID.String(), uuid.New(), ext)
}
