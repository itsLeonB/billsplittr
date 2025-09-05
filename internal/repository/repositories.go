package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/entity"
	crud "github.com/itsLeonB/go-crud"
)

type GroupExpenseRepository interface {
	crud.CRUDRepository[entity.GroupExpense]
	SyncParticipants(ctx context.Context, groupExpenseID uuid.UUID, participants []entity.ExpenseParticipant) error
}

type ExpenseItemRepository interface {
	crud.CRUDRepository[entity.ExpenseItem]
	SyncParticipants(ctx context.Context, expenseItemID uuid.UUID, participants []entity.ItemParticipant) error
}

type ExpenseParticipantRepository interface {
	crud.CRUDRepository[entity.ExpenseParticipant]
}

type OtherFeeRepository interface {
	crud.CRUDRepository[entity.OtherFee]
	SyncParticipants(ctx context.Context, feeID uuid.UUID, participants []entity.FeeParticipant) error
}

type ExpenseBillRepository interface {
	crud.CRUDRepository[entity.ExpenseBill]
}

// StorageRepository handles file storage operations
type StorageRepository interface {
	Upload(ctx context.Context, req *entity.StorageUploadRequest) (*entity.StorageUploadResponse, error)
	Download(ctx context.Context, bucketName, objectKey string) ([]byte, error)
	Delete(ctx context.Context, bucketName, objectKey string) error
	GetSignedURL(ctx context.Context, bucketName, objectKey string, expiration time.Duration) (string, error)
	Close() error
}
