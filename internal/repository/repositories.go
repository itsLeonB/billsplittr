package repository

import (
	"context"

	"github.com/itsLeonB/billsplittr/internal/entity"
)

type UserRepository interface {
	Insert(ctx context.Context, user entity.User) (entity.User, error)
	Find(ctx context.Context, spec entity.User) (entity.User, error)
}
