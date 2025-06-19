package mapper

import (
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
)

func UserToAuthData(user entity.User) map[string]any {
	return map[string]any{
		appconstant.ContextUserID: user.ID,
	}
}

func UserToResponse(user entity.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}

func UserToProfileResponse(user entity.User) dto.ProfileResponse {
	return dto.ProfileResponse{
		UserID:    user.ID,
		ProfileID: user.Profile.ID,
		Name:      user.Profile.Name,
		CreatedAt: user.Profile.CreatedAt,
		UpdatedAt: user.Profile.UpdatedAt,
		DeletedAt: user.Profile.DeletedAt,
	}
}
