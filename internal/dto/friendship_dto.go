package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
)

type NewAnonymousFriendshipRequest struct {
	UserID uuid.UUID
	Name   string `json:"name" binding:"required,min=3"`
}

type FriendshipResponse struct {
	ID          uuid.UUID                  `json:"id"`
	Type        appconstant.FriendshipType `json:"type"`
	ProfileID   uuid.UUID                  `json:"profileId"`
	ProfileName string                     `json:"profileName"`
	CreatedAt   time.Time                  `json:"createdAt"`
	UpdatedAt   time.Time                  `json:"updatedAt"`
	DeletedAt   time.Time                  `json:"deletedAt,omitzero"`
}
