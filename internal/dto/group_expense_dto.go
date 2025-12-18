package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/shopspring/decimal"
)

type NewGroupExpenseRequest struct {
	CreatorProfileID uuid.UUID       `validate:"required"`
	PayerProfileID   uuid.UUID       `validate:"required"`
	TotalAmount      decimal.Decimal `validate:"required"`
	// Deprecated: use ItemsTotal instead
	Subtotal    decimal.Decimal `validate:"required"`
	ItemsTotal  decimal.Decimal
	FeesTotal   decimal.Decimal
	Description string
	Items       []ExpenseItemData `validate:"required,min=1,dive"`
	OtherFees   []OtherFeeData    `validate:"dive"`
}

type NewDraftExpense struct {
	CreatorProfileID uuid.UUID `validate:"required"`
	Description      string
}

type GroupExpenseResponse struct {
	ID             uuid.UUID
	PayerProfileID uuid.UUID
	TotalAmount    decimal.Decimal
	// Deprecated: use ItemsTotal instead
	Subtotal         decimal.Decimal
	ItemsTotal       decimal.Decimal
	FeesTotal        decimal.Decimal
	Description      string
	CreatorProfileID uuid.UUID
	// Deprecated: refer to Status instead
	Confirmed bool
	// Deprecated: refer to Status instead
	ParticipantsConfirmed bool
	Status                appconstant.ExpenseStatus
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             time.Time

	// Relationships
	Items        []ExpenseItemResponse
	OtherFees    []OtherFeeResponse
	Participants []ExpenseParticipantResponse
	Bill         ExpenseBillResponse
}

type ExpenseParticipantResponse struct {
	ProfileID   uuid.UUID
	ShareAmount decimal.Decimal
}
