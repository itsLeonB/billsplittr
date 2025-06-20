package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
)

type AuthService interface {
	Register(ctx context.Context, request dto.RegisterRequest) error
	Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error)
}

type UserService interface {
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
	GetByID(ctx context.Context, id uuid.UUID) (dto.UserResponse, error)
	GetProfile(ctx context.Context, id uuid.UUID) (dto.ProfileResponse, error)
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (entity.User, error)
}

type FriendshipService interface {
	CreateAnonymous(ctx context.Context, request dto.NewAnonymousFriendshipRequest) (dto.FriendshipResponse, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]dto.FriendshipResponse, error)
}

type DebtService interface {
	RecordNewTransaction(ctx context.Context, request dto.NewDebtTransactionRequest) (dto.DebtTransactionResponse, error)
}

type TransferMethodService interface {
	GetAll(ctx context.Context) ([]dto.TransferMethodResponse, error)
}
