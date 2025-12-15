package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type NewExpenseItemRequest struct {
	ProfileID      uuid.UUID `validate:"required"`
	GroupExpenseID uuid.UUID `validate:"required"`
	ExpenseItemData
}

type UpdateExpenseItemRequest struct {
	ProfileID      uuid.UUID `validate:"required"`
	ID             uuid.UUID `validate:"required"`
	GroupExpenseID uuid.UUID `validate:"required"`
	ExpenseItemData
	Participants []ItemParticipantData `validate:"dive"`
}

type ExpenseItemResponse struct {
	ID             uuid.UUID
	GroupExpenseID uuid.UUID
	Name           string
	Amount         decimal.Decimal
	Quantity       int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
	Participants   []ItemParticipantData
}

type ItemParticipantData struct {
	ProfileID uuid.UUID       `validate:"required"`
	Share     decimal.Decimal `validate:"required"`
}

type ExpenseItemData struct {
	Name     string          `validate:"required,min=3"`
	Amount   decimal.Decimal `validate:"required"`
	Quantity int             `validate:"required,min=1"`
}
