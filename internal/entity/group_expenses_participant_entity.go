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

func (gep GroupExpenseParticipant) IsZero() bool {
	return gep.ID == uuid.Nil
}
