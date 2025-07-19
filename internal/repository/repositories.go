package repository

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/ezutil"
)

type UserRepository interface {
	Insert(ctx context.Context, user entity.User) (entity.User, error)
	Find(ctx context.Context, spec entity.User) (entity.User, error)
}

type UserProfileRepository interface {
	Insert(ctx context.Context, profile entity.UserProfile) (entity.UserProfile, error)
}

type FriendshipRepository interface {
	Insert(ctx context.Context, friendship entity.Friendship) (entity.Friendship, error)
	FindAll(ctx context.Context, spec entity.FriendshipSpecification) ([]entity.Friendship, error)
	FindFirst(ctx context.Context, spec entity.FriendshipSpecification) (entity.Friendship, error)
	FindByProfileIDs(ctx context.Context, profileID1, profileID2 uuid.UUID) (entity.Friendship, error)
}

type DebtTransactionRepository interface {
	ezutil.CRUDRepository[entity.DebtTransaction]
	FindAllByProfileID(ctx context.Context, userProfileID, friendProfileID uuid.UUID) ([]entity.DebtTransaction, error)
	FindAllByUserProfileID(ctx context.Context, userProfileID uuid.UUID) ([]entity.DebtTransaction, error)
}

type TransferMethodRepository interface {
	ezutil.CRUDRepository[entity.TransferMethod]
}

type GroupExpenseRepository interface {
	ezutil.CRUDRepository[entity.GroupExpense]
	SyncParticipants(ctx context.Context, groupExpenseID uuid.UUID, participants []entity.ExpenseParticipant) error
}

type ExpenseItemRepository interface {
	ezutil.CRUDRepository[entity.ExpenseItem]
	SyncParticipants(ctx context.Context, expenseItemID uuid.UUID, participants []entity.ItemParticipant) error
}

type ExpenseParticipantRepository interface {
	ezutil.CRUDRepository[entity.ExpenseParticipant]
}

type OtherFeeRepository interface {
	ezutil.CRUDRepository[entity.OtherFee]
	SyncParticipants(ctx context.Context, feeID uuid.UUID, participants []entity.FeeParticipant) error
}

// ImageRepository defines the behavior for image storage.
type ImageRepository interface {
	Upload(ctx context.Context, reader io.Reader, contentType string) (string, error)
	GenerateSignedURL(ctx context.Context, objectName string, duration time.Duration) (string, error)
	Delete(ctx context.Context, objectName string) error
}

type ExpenseBillRepository interface {
	ezutil.CRUDRepository[entity.ExpenseBill]
}
