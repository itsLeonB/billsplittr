package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/mapper"
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/itsLeonB/cocoon-protos/gen/go/friendship/v1"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/rotisserie/eris"
)

type friendshipServiceImpl struct {
	debtTransactionRepository repository.DebtTransactionRepository
	friendshipClient          friendship.FriendshipServiceClient
}

func NewFriendshipService(
	debtTransactionRepository repository.DebtTransactionRepository,
	friendshipClient friendship.FriendshipServiceClient,
) FriendshipService {
	return &friendshipServiceImpl{
		debtTransactionRepository,
		friendshipClient,
	}
}

func (fs *friendshipServiceImpl) CreateAnonymous(
	ctx context.Context,
	request dto.NewAnonymousFriendshipRequest,
) (dto.FriendshipResponse, error) {
	req := &friendship.CreateAnonymousRequest{
		ProfileId: request.ProfileID.String(),
		Name:      request.Name,
	}

	response, err := fs.friendshipClient.CreateAnonymous(ctx, req)
	if err != nil {
		return dto.FriendshipResponse{}, eris.Wrap(err, appconstant.ErrServiceClient)
	}

	return mapper.FromFriendshipResponseProto(response.GetFriendship())
}

func (fs *friendshipServiceImpl) GetAll(ctx context.Context, profileID uuid.UUID) ([]dto.FriendshipResponse, error) {
	response, err := fs.friendshipClient.GetAll(ctx, &friendship.GetAllRequest{ProfileId: profileID.String()})
	if err != nil {
		return nil, eris.Wrap(err, appconstant.ErrServiceClient)
	}

	return ezutil.MapSliceWithError(response.GetFriendships(), mapper.FromFriendshipResponseProto)
}

func (fs *friendshipServiceImpl) GetDetails(ctx context.Context, profileID, friendshipID uuid.UUID) (dto.FriendDetailsResponse, error) {
	req := &friendship.GetDetailsRequest{
		ProfileId:    profileID.String(),
		FriendshipId: friendshipID.String(),
	}

	response, err := fs.friendshipClient.GetDetails(ctx, req)
	if err != nil {
		return dto.FriendDetailsResponse{}, eris.Wrap(err, appconstant.ErrServiceClient)
	}

	profileID1, err := ezutil.Parse[uuid.UUID](response.GetProfileId1())
	if err != nil {
		return dto.FriendDetailsResponse{}, err
	}

	profileID2, err := ezutil.Parse[uuid.UUID](response.GetProfileId2())
	if err != nil {
		return dto.FriendDetailsResponse{}, err
	}

	debtTransactions, err := fs.debtTransactionRepository.FindAllByProfileID(ctx, profileID1, profileID2)
	if err != nil {
		return dto.FriendDetailsResponse{}, err
	}

	return mapper.MapToFriendDetailsResponse(profileID, response, debtTransactions)
}

func (fs *friendshipServiceImpl) IsFriends(ctx context.Context, profileID1, profileID2 uuid.UUID) (bool, bool, error) {
	req := friendship.IsFriendsRequest{
		ProfileId_1: profileID1.String(),
		ProfileId_2: profileID2.String(),
	}

	response, err := fs.friendshipClient.IsFriends(ctx, &req)
	if err != nil {
		return false, false, eris.Wrap(err, appconstant.ErrServiceClient)
	}

	return response.GetIsFriends(), response.GetIsAnonymous(), nil
}
