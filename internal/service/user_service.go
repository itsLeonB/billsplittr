package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/mapper"
	"github.com/itsLeonB/billsplittr/internal/repository"
)

type userServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userServiceImpl{
		userRepository: userRepository,
	}
}

func (us *userServiceImpl) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	spec := entity.User{}
	spec.ID = id

	user, err := us.userRepository.Find(ctx, spec)
	if err != nil {
		return false, err
	}
	if user.IsZero() {
		return false, nil
	}

	return true, nil
}

func (us *userServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (dto.UserResponse, error) {
	spec := entity.User{}
	spec.ID = id

	user, err := us.userRepository.Find(ctx, spec)
	if err != nil {
		return dto.UserResponse{}, err
	}
	if user.IsZero() {
		return dto.UserResponse{}, nil
	}

	return mapper.UserToResponse(user), nil
}
