package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type NewGroupExpenseRequest struct {
	CreatorProfileID uuid.UUID       `validate:"required"`
	PayerProfileID   uuid.UUID       `validate:"required"`
	TotalAmount      decimal.Decimal `validate:"required"`
	Subtotal         decimal.Decimal `validate:"required"`
	Description      string
	Items            []ExpenseItemData `validate:"required,min=1,dive"`
	OtherFees        []OtherFeeData    `validate:"required,dive"`
}

type GroupExpenseResponse struct {
	ID                    uuid.UUID
	PayerProfileID        uuid.UUID
	PayerName             string
	PaidByUser            bool
	TotalAmount           decimal.Decimal
	Description           string
	Items                 []ExpenseItemResponse
	OtherFees             []OtherFeeResponse
	CreatorProfileID      uuid.UUID
	CreatorName           string
	CreatedByUser         bool
	Confirmed             bool
	ParticipantsConfirmed bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             time.Time
	Participants          []ExpenseParticipantResponse
}

type ExpenseParticipantResponse struct {
	ProfileID   uuid.UUID
	ShareAmount decimal.Decimal
}
