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
