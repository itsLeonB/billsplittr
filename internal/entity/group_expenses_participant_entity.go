package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type GroupExpenseParticipant struct {
	BaseEntity
	GroupExpenseID       uuid.UUID
	ParticipantProfileID uuid.UUID
	ShareAmount          decimal.Decimal
	Description          string
}
