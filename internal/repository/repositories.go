package repository

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/entity"
	crud "github.com/itsLeonB/go-crud"
)

type DebtTransactionRepository interface {
	crud.CRUDRepository[entity.DebtTransaction]
	FindAllByProfileID(ctx context.Context, userProfileID, friendProfileID uuid.UUID) ([]entity.DebtTransaction, error)
	FindAllByUserProfileID(ctx context.Context, userProfileID uuid.UUID) ([]entity.DebtTransaction, error)
}

type TransferMethodRepository interface {
	crud.CRUDRepository[entity.TransferMethod]
}

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

// ImageRepository defines the behavior for image storage.
type ImageRepository interface {
	Upload(ctx context.Context, reader io.Reader, contentType string) (string, error)
	GenerateSignedURL(ctx context.Context, objectName string, duration time.Duration) (string, error)
	Delete(ctx context.Context, objectName string) error
}

type ExpenseBillRepository interface {
	crud.CRUDRepository[entity.ExpenseBill]
}
