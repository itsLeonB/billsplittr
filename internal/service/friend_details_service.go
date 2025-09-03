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

type friendDetailsServiceImpl struct {
	friendshipClient friendship.FriendshipServiceClient
	debtSvc          DebtService
}

func NewFriendDetailsService(friendshipClient friendship.FriendshipServiceClient, debtSvc DebtService) FriendDetailsService {
	return &friendDetailsServiceImpl{friendshipClient, debtSvc}
}

func (fds *friendDetailsServiceImpl) GetDetails(ctx context.Context, profileID, friendshipID uuid.UUID) (dto.FriendDetailsResponse, error) {
	req := &friendship.GetDetailsRequest{
		ProfileId:    profileID.String(),
		FriendshipId: friendshipID.String(),
	}

	response, err := fds.friendshipClient.GetDetails(ctx, req)
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

	friendProfileID := profileID2
	if ezutil.CompareUUID(profileID, profileID2) == 0 {
		friendProfileID = profileID1
	}

	debtTransactions, err := fds.debtSvc.GetAllByProfileIDs(ctx, profileID, friendProfileID)
	if err != nil {
		return dto.FriendDetailsResponse{}, err
	}

	return mapper.MapToFriendDetailsResponse(profileID, response, debtTransactions)
}
