package dto

import (
	"time"

	"github.com/google/uuid"
)

type NewExpenseBillRequest struct {
	CreatorProfileID uuid.UUID `validate:"required"`
	PayerProfileID   uuid.UUID `validate:"required"`
	Filename         string    `validate:"required,min=3"`
}

type ExpenseBillResponse struct {
	ID               uuid.UUID
	CreatorProfileID uuid.UUID
	PayerProfileID   uuid.UUID
	GroupExpenseID   uuid.UUID
	Filename         string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time
}
