package mapper

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/util"
	"github.com/rotisserie/eris"
)

func FriendshipToResponse(userProfileID uuid.UUID, friendship entity.Friendship) (dto.FriendshipResponse, error) {
	friendProfile, err := selectFriendProfile(userProfileID, friendship)
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

func selectFriendProfile(userProfileID uuid.UUID, friendship entity.Friendship) (entity.UserProfile, error) {
	switch userProfileID {
	case friendship.ProfileID1:
		return friendship.Profile2, nil
	case friendship.ProfileID2:
		return friendship.Profile1, nil
	default:
		return entity.UserProfile{}, eris.New(fmt.Sprintf(
			"mismatched user profile ID: %s with friendship ID: %s",
			userProfileID,
			friendship.ID,
		))
	}
}

func OrderProfilesToFriendship(userProfile, friendProfile entity.UserProfile) (entity.Friendship, error) {
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
