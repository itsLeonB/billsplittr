package repository

import (
	"context"

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
