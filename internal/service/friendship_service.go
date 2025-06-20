package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/mapper"
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/itsLeonB/billsplittr/internal/util"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
)

type friendshipServiceImpl struct {
	transactor            ezutil.Transactor
	userProfileRepository repository.UserProfileRepository
	friendshipRepository  repository.FriendshipRepository
	userService           UserService
}

func NewFriendshipRepository(
	transactor ezutil.Transactor,
	userProfileRepository repository.UserProfileRepository,
	friendshipRepository repository.FriendshipRepository,
	userService UserService,
) FriendshipService {
	return &friendshipServiceImpl{
		transactor,
		userProfileRepository,
		friendshipRepository,
		userService,
	}
}

func (fs *friendshipServiceImpl) CreateAnonymous(
	ctx context.Context,
	request dto.NewAnonymousFriendshipRequest,
) (dto.FriendshipResponse, error) {
	var response dto.FriendshipResponse

	err := fs.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		user, err := fs.userService.GetByIDForUpdate(ctx, request.UserID)
		if err != nil {
			return err
		}

		if err = fs.validateExistingAnonymousFriendship(ctx, user.Profile.ID, request.Name); err != nil {
			return err
		}

		response, err = fs.insertAnonymousFriendship(ctx, user.Profile, request.Name)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return dto.FriendshipResponse{}, err
	}

	return response, nil
}

func (fs *friendshipServiceImpl) GetAll(ctx context.Context, userID uuid.UUID) ([]dto.FriendshipResponse, error) {
	user, err := fs.userService.GetByIDForUpdate(ctx, userID)
	if err != nil {
		return nil, err
	}
	spec := entity.FriendshipSpecification{ProfileID1: user.Profile.ID}
	spec.PreloadRelations = []string{"Profile1", "Profile2"}

	friendships, err := fs.friendshipRepository.FindAll(ctx, spec)
	if err != nil {
		return nil, err
	}

	mapperFunc := func(friendship entity.Friendship) (dto.FriendshipResponse, error) {
		return mapper.FriendshipToResponse(user.Profile.ID, friendship)
	}

	return ezutil.MapSliceWithError(friendships, mapperFunc)
}

func (fs *friendshipServiceImpl) validateExistingAnonymousFriendship(
	ctx context.Context,
	userProfileID uuid.UUID,
	friendName string,
) error {
	friendshipSpec := entity.FriendshipSpecification{}
	friendshipSpec.ProfileID1 = userProfileID
	friendshipSpec.Name = friendName
	friendshipSpec.Type = appconstant.Anonymous

	existingFriendship, err := fs.friendshipRepository.FindFirst(ctx, friendshipSpec)
	if err != nil {
		return err
	}
	if !existingFriendship.IsZero() && !existingFriendship.IsDeleted() {
		return ezutil.ConflictError(fmt.Sprintf("anonymous friend named %s already exists", friendName))
	}

	return nil
}

func (fs *friendshipServiceImpl) insertAnonymousFriendship(
	ctx context.Context,
	userProfile entity.UserProfile,
	friendName string,
) (dto.FriendshipResponse, error) {
	newProfile := entity.UserProfile{Name: friendName}

	insertedProfile, err := fs.userProfileRepository.Insert(ctx, newProfile)
	if err != nil {
		return dto.FriendshipResponse{}, err
	}

	newFriendship, err := orderProfileID(userProfile, insertedProfile)
	if err != nil {
		return dto.FriendshipResponse{}, err
	}

	newFriendship.Type = appconstant.Anonymous

	insertedFriendship, err := fs.friendshipRepository.Insert(ctx, newFriendship)
	if err != nil {
		return dto.FriendshipResponse{}, err
	}

	return mapper.FriendshipToResponse(userProfile.ID, insertedFriendship)
}

func orderProfileID(userProfile, friendProfile entity.UserProfile) (entity.Friendship, error) {
	switch util.CompareUUID(userProfile.ID, friendProfile.ID) {
	case 1:
		return entity.Friendship{
			ProfileID1: friendProfile.ID,
			ProfileID2: userProfile.ID,
			Profile1:   friendProfile,
			Profile2:   userProfile,
		}, nil
	case -1:
		return entity.Friendship{
			ProfileID1: userProfile.ID,
			ProfileID2: friendProfile.ID,
			Profile1:   userProfile,
			Profile2:   friendProfile,
		}, nil
	default:
		return entity.Friendship{}, eris.New("both IDs are equal, cannot create friendship")
	}
}
