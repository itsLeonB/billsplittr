package mapper

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/cocoon-protos/gen/go/profile/v1"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
)

func FromProfileProto(res *profile.ProfileResponse) (dto.ProfileResponse, error) {
	id, err := ezutil.Parse[uuid.UUID](res.GetId())
	if err != nil {
		return dto.ProfileResponse{}, err
	}

	userID := uuid.Nil
	if res.GetUserId() != "" {
		userID, err = ezutil.Parse[uuid.UUID](res.GetUserId())
		if err != nil {
			return dto.ProfileResponse{}, err
		}
	}

	if res.GetIsAnonymous() && userID != uuid.Nil {
		return dto.ProfileResponse{}, eris.Errorf("anonymous user has userID: %s", userID)
	}

	return dto.ProfileResponse{
		ID:          id,
		UserID:      userID,
		Name:        res.GetName(),
		CreatedAt:   res.GetCreatedAt().AsTime(),
		UpdatedAt:   res.GetUpdatedAt().AsTime(),
		DeletedAt:   res.GetDeletedAt().AsTime(),
		IsAnonymous: res.GetIsAnonymous(),
	}, nil
}
