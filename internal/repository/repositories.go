package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/entity"
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
}

type DebtTransactionRepository interface {
	Insert(ctx context.Context, debtTransaction entity.DebtTransaction) (entity.DebtTransaction, error)
	FindAllByProfileID(ctx context.Context, userProfileID, friendProfileID uuid.UUID) ([]entity.DebtTransaction, error)
	FindAllByUserProfileID(ctx context.Context, userProfileID uuid.UUID) ([]entity.DebtTransaction, error)
}

type TransferMethodRepository interface {
	FindAll(ctx context.Context, spec entity.TransferMethod) ([]entity.TransferMethod, error)
	FindFirst(ctx context.Context, spec entity.TransferMethod) (entity.TransferMethod, error)
}
