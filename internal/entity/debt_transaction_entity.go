package entity

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/shopspring/decimal"
)

type DebtTransaction struct {
	BaseEntity
	LenderProfileID   uuid.UUID
	BorrowerProfileID uuid.UUID
	Type              appconstant.DebtTransactionType
	Amount            decimal.Decimal
	TransferMethodID  uuid.UUID
	Description       string
}

func (dt DebtTransaction) IsZero() bool {
	return dt.ID == uuid.Nil
}
