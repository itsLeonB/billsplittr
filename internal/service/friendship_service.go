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
	userRepository        repository.UserRepository
	userProfileRepository repository.UserProfileRepository
	friendshipRepository  repository.FriendshipRepository
}

func NewFriendshipRepository(
	transactor ezutil.Transactor,
	userRepository repository.UserRepository,
	userProfileRepository repository.UserProfileRepository,
	friendshipRepository repository.FriendshipRepository,
) FriendshipService {
	return &friendshipServiceImpl{
		transactor,
		userRepository,
		userProfileRepository,
		friendshipRepository,
	}
}

func (fs *friendshipServiceImpl) CreateAnonymous(
	ctx context.Context,
	request dto.NewAnonymousFriendshipRequest,
) (dto.FriendshipResponse, error) {
	var response dto.FriendshipResponse

	err := ezutil.WithinTransaction(ctx, fs.transactor, func(ctx context.Context) error {
		user, err := fs.findUser(ctx, request.UserID)
		if err != nil {
			return err
		}

		userProfileID := user.Profile.ID

		if err = fs.validateExistingAnonymousFriendship(ctx, userProfileID, request.Name); err != nil {
			return err
		}

		response, err = fs.insertAnonymousFriendship(ctx, userProfileID, request.Name)
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

func (fs *friendshipServiceImpl) findUser(ctx context.Context, id uuid.UUID) (entity.User, error) {
	userSpec := entity.User{}
	userSpec.ID = id

	user, err := fs.userRepository.Find(ctx, userSpec)
	if err != nil {
		return entity.User{}, err
	}
	if user.IsZero() {
		return entity.User{}, ezutil.NotFoundError(fmt.Sprintf(appconstant.MsgUserNotFound, id))
	}
	if user.IsDeleted() {
		return entity.User{}, ezutil.UnprocessableEntityError(fmt.Sprintf(appconstant.MsgUserDeleted, id))
	}

	return user, nil
}

func (fs *friendshipServiceImpl) validateExistingAnonymousFriendship(
	ctx context.Context,
	userProfileID uuid.UUID,
	friendName string,
) error {
	friendshipSpec := entity.FriendshipSpecification{}
	friendshipSpec.ProfileID = userProfileID
	friendshipSpec.Name = friendName
	friendshipSpec.Type = appconstant.Anonymous

	existingFriendship, err := fs.friendshipRepository.FindFirst(ctx, friendshipSpec)
	if err != nil {
		return err
	}
	if !existingFriendship.IsZero() && !existingFriendship.IsDeleted() {
		return ezutil.ConflictError(fmt.Sprintf("anonymous friend named: %s already exists", friendName))
	}

	return nil
}

func (fs *friendshipServiceImpl) insertAnonymousFriendship(
	ctx context.Context,
	userProfileID uuid.UUID,
	friendName string,
) (dto.FriendshipResponse, error) {
	newProfile := entity.UserProfile{Name: friendName}

	insertedProfile, err := fs.userProfileRepository.Insert(ctx, newProfile)
	if err != nil {
		return dto.FriendshipResponse{}, err
	}

	newFriendship, err := orderProfileID(userProfileID, insertedProfile.ID)
	if err != nil {
		return dto.FriendshipResponse{}, err
	}

	newFriendship.Type = appconstant.Anonymous

	insertedFriendship, err := fs.friendshipRepository.Insert(ctx, newFriendship)
	if err != nil {
		return dto.FriendshipResponse{}, err
	}

	return mapper.FriendshipToResponse(userProfileID, insertedFriendship)
}

func orderProfileID(profileID1, profileID2 uuid.UUID) (entity.Friendship, error) {
	switch util.CompareUUID(profileID1, profileID2) {
	case 1:
		return entity.Friendship{
			ProfileID1: profileID2,
			ProfileID2: profileID1,
		}, nil
	case -1:
		return entity.Friendship{
			ProfileID1: profileID1,
			ProfileID2: profileID2,
		}, nil
	default:
		return entity.Friendship{}, eris.New("both IDs are equal, cannot create friendship")
	}
}
