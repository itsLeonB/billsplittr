package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type GroupExpenses struct {
	BaseEntity
	PayerProfileID uuid.UUID
	TotalAmount    decimal.Decimal
	Description    string
}

func (ge GroupExpenses) IsZero() bool {
	return ge.ID == uuid.Nil
}
