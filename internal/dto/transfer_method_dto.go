package dto

import (
	"time"

	"github.com/google/uuid"
)

type TransferMethodResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Display   string    `json:"display"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt,omitzero"`
}
