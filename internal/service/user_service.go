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
	userRepository        repository.UserRepository
	userProfileRepository repository.UserProfileRepository
}

func NewUserService(
	userRepository repository.UserRepository,
	userProfileRepository repository.UserProfileRepository,
) UserService {
	return &userServiceImpl{
		userRepository,
		userProfileRepository,
	}
}

func (us *userServiceImpl) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	user, err := us.findById(ctx, id)
	if err != nil {
		return false, err
	}

	return !user.IsZero(), nil
}

func (us *userServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (dto.UserResponse, error) {
	user, err := us.findById(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return mapper.UserToResponse(user), nil
}

func (us *userServiceImpl) GetProfile(ctx context.Context, id uuid.UUID) (dto.ProfileResponse, error) {
	user, err := us.findById(ctx, id)
	if err != nil {
		return dto.ProfileResponse{}, err
	}

	return mapper.UserToProfileResponse(user), nil
}

func (us *userServiceImpl) findById(ctx context.Context, id uuid.UUID) (entity.User, error) {
	spec := entity.User{}
	spec.ID = id

	user, err := us.userRepository.Find(ctx, spec)
	if err != nil {
		return entity.User{}, err
	}
	if user.IsZero() {
		return entity.User{}, nil
	}

	return user, nil
}
