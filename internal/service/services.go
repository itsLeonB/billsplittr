package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
)

type AuthService interface {
	Register(ctx context.Context, request dto.RegisterRequest) error
	Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error)
}

type UserService interface {
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
	GetByID(ctx context.Context, id uuid.UUID) (dto.UserResponse, error)
	GetProfile(ctx context.Context, id uuid.UUID) (dto.ProfileResponse, error)
}

type FriendshipService interface {
	CreateAnonymous(ctx context.Context, request dto.NewAnonymousFriendshipRequest) (dto.FriendshipResponse, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]dto.FriendshipResponse, error)
}
