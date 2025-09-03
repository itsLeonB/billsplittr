package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/mapper"
	"github.com/itsLeonB/cocoon-protos/gen/go/friendship/v1"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/rotisserie/eris"
)

type friendshipServiceImpl struct {
	friendshipClient friendship.FriendshipServiceClient
}

func NewFriendshipService(friendshipClient friendship.FriendshipServiceClient) FriendshipService {
	return &friendshipServiceImpl{friendshipClient}
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
