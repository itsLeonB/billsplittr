package mapper

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
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
