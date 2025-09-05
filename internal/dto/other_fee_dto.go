package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/shopspring/decimal"
)

type NewOtherFeeRequest struct {
	ProfileID      uuid.UUID `validate:"required"`
	GroupExpenseID uuid.UUID `validate:"required"`
	OtherFeeData   `validate:"required,dive"`
}

type UpdateOtherFeeRequest struct {
	ProfileID      uuid.UUID `validate:"required"`
	ID             uuid.UUID `validate:"required"`
	GroupExpenseID uuid.UUID `validate:"required"`
	OtherFeeData   `validate:"required,dive"`
}

type OtherFeeResponse struct {
	ID                uuid.UUID
	GroupExpenseID    uuid.UUID
	Name              string
	Amount            decimal.Decimal
	CalculationMethod appconstant.FeeCalculationMethod
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         time.Time
	Participants      []FeeParticipantResponse
}

type FeeParticipantResponse struct {
	ProfileID   uuid.UUID
	ShareAmount decimal.Decimal
}

type FeeCalculationMethodInfo struct {
	Method      appconstant.FeeCalculationMethod
	Display     string
	Description string
}

type OtherFeeData struct {
	Name              string                           `validate:"required,min=3"`
	Amount            decimal.Decimal                  `validate:"required"`
	CalculationMethod appconstant.FeeCalculationMethod `validate:"required"`
}
