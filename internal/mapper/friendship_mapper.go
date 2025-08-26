package mapper

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/cocoon-protos/gen/go/friendship"
	"github.com/itsLeonB/ezutil"
)

func MapToFriendDetailsResponse(
	userProfileID uuid.UUID,
	friendDetails *friendship.FriendDetails,
	debtTransactions []entity.DebtTransaction,
) (dto.FriendDetailsResponse, error) {
	id, err := ezutil.Parse[uuid.UUID](friendDetails.GetId())
	if err != nil {
		return dto.FriendDetailsResponse{}, err
	}

	profileID, err := ezutil.Parse[uuid.UUID](friendDetails.GetProfileId())
	if err != nil {
		return dto.FriendDetailsResponse{}, err
	}

	return dto.FriendDetailsResponse{
		Friend: dto.FriendDetails{
			ID:        id,
			ProfileID: profileID,
			Name:      friendDetails.GetName(),
			Type:      appconstant.FriendshipType(friendDetails.GetType().String()),
			CreatedAt: friendDetails.GetCreatedAt().AsTime(),
			UpdatedAt: friendDetails.GetUpdatedAt().AsTime(),
			DeletedAt: friendDetails.GetDeletedAt().AsTime(),
		},
		Balance:      MapToFriendBalanceSummary(userProfileID, debtTransactions),
		Transactions: ezutil.MapSlice(debtTransactions, GetDebtTransactionSimpleMapper(userProfileID)),
	}, nil
}

func FromFriendshipResponseProto(response *friendship.FriendshipResponse) (dto.FriendshipResponse, error) {
	id, err := ezutil.Parse[uuid.UUID](response.GetId())
	if err != nil {
		return dto.FriendshipResponse{}, err
	}

	profileID, err := ezutil.Parse[uuid.UUID](response.GetProfileId())
	if err != nil {
		return dto.FriendshipResponse{}, err
	}

	return dto.FriendshipResponse{
		ID:          id,
		Type:        appconstant.FriendshipType(response.GetType().String()),
		ProfileID:   profileID,
		ProfileName: response.GetProfileName(),
		CreatedAt:   response.GetCreatedAt().AsTime(),
		UpdatedAt:   response.GetUpdatedAt().AsTime(),
		DeletedAt:   response.GetDeletedAt().AsTime(),
	}, nil
}
