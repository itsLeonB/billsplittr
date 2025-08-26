package mapper

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/cocoon-protos/gen/go/friendship/v1"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
)

func MapToFriendDetailsResponse(
	userProfileID uuid.UUID,
	friendDetails *friendship.GetDetailsResponse,
	debtTransactions []entity.DebtTransaction,
) (dto.FriendDetailsResponse, error) {
	if friendDetails == nil {
		return dto.FriendDetailsResponse{}, eris.New("friendDetails is nil")
	}

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
			Type:      FromProtoFriendshipType(friendDetails.Type),
			CreatedAt: FromProtoTime(friendDetails.GetCreatedAt()),
			UpdatedAt: FromProtoTime(friendDetails.GetUpdatedAt()),
			DeletedAt: FromProtoTime(friendDetails.GetDeletedAt()),
		},
		Balance:      MapToFriendBalanceSummary(userProfileID, debtTransactions),
		Transactions: ezutil.MapSlice(debtTransactions, GetDebtTransactionSimpleMapper(userProfileID)),
	}, nil
}

func FromFriendshipResponseProto(response *friendship.FriendshipResponse) (dto.FriendshipResponse, error) {
	if response == nil {
		return dto.FriendshipResponse{}, eris.New("proto is nil")
	}

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
		Type:        FromProtoFriendshipType(response.GetType()),
		ProfileID:   profileID,
		ProfileName: response.GetProfileName(),
		CreatedAt:   FromProtoTime(response.GetCreatedAt()),
		UpdatedAt:   FromProtoTime(response.GetUpdatedAt()),
		DeletedAt:   FromProtoTime(response.GetDeletedAt()),
	}, nil
}

func FromProtoFriendshipType(ft friendship.FriendshipType) appconstant.FriendshipType {
	switch ft {
	case friendship.FriendshipType_FRIENDSHIP_TYPE_REAL:
		return appconstant.Real
	case friendship.FriendshipType_FRIENDSHIP_TYPE_ANON:
		return appconstant.Anonymous
	default:
		return ""
	}
}
