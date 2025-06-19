package repository

import (
	"context"

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
