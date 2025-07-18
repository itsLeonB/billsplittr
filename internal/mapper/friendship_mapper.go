package mapper

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/helper"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
)

func FriendshipToResponse(userProfileID uuid.UUID, friendship entity.Friendship) (dto.FriendshipResponse, error) {
	_, friendProfile, err := helper.SelectProfiles(userProfileID, friendship)
	if err != nil {
		return dto.FriendshipResponse{}, err
	}

	return dto.FriendshipResponse{
		ID:          friendship.ID,
		Type:        friendship.Type,
		ProfileID:   friendProfile.ID,
		ProfileName: friendProfile.Name,
		CreatedAt:   friendship.CreatedAt,
		UpdatedAt:   friendship.UpdatedAt,
		DeletedAt:   friendship.DeletedAt.Time,
	}, nil
}

func OrderProfilesToFriendship(userProfile, friendProfile entity.UserProfile) (entity.Friendship, error) {
	switch ezutil.CompareUUID(userProfile.ID, friendProfile.ID) {
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

func MapToFriendshipWithProfile(userProfileID uuid.UUID, friendship entity.Friendship) (dto.FriendshipWithProfile, error) {
	friendshipResponse, err := FriendshipToResponse(userProfileID, friendship)
	if err != nil {
		return dto.FriendshipWithProfile{}, err
	}

	userProfile, friendProfile, err := helper.SelectProfiles(userProfileID, friendship)
	if err != nil {
		return dto.FriendshipWithProfile{}, err
	}

	return dto.FriendshipWithProfile{
		Friendship:    friendshipResponse,
		UserProfile:   ProfileToResponse(userProfile),
		FriendProfile: ProfileToResponse(friendProfile),
	}, nil
}

func MapToFriendDetailsResponse(
	userProfileID uuid.UUID,
	friendship entity.Friendship,
	debtTransactions []entity.DebtTransaction,
) (dto.FriendDetailsResponse, error) {
	friendshipWithProfile, err := MapToFriendshipWithProfile(userProfileID, friendship)
	if err != nil {
		return dto.FriendDetailsResponse{}, err
	}

	friendProfile := friendshipWithProfile.FriendProfile

	return dto.FriendDetailsResponse{
		Friend: dto.FriendDetails{
			ID:        friendship.ID,
			ProfileID: friendProfile.ProfileID,
			Name:      friendProfile.Name,
			Type:      friendship.Type,
			CreatedAt: friendship.CreatedAt,
			UpdatedAt: friendship.UpdatedAt,
			DeletedAt: friendship.DeletedAt.Time,
		},
		Balance:      MapToFriendBalanceSummary(userProfileID, debtTransactions),
		Transactions: ezutil.MapSlice(debtTransactions, GetDebtTransactionSimpleMapper(userProfileID)),
	}, nil
}
