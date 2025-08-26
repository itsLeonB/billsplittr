package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ExpenseParticipant struct {
	BaseEntity
	GroupExpenseID       uuid.UUID
	ParticipantProfileID uuid.UUID
	ShareAmount          decimal.Decimal
	Description          string
	Confirmed            bool
}

func (ep ExpenseParticipant) TableName() string {
	return "group_expense_participants"
}
