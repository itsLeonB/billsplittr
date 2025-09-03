package entity

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/go-crud"
)

type ExpenseBill struct {
	crud.BaseEntity
	PayerProfileID   uuid.UUID
	ImageName        string
	GroupExpenseID   uuid.NullUUID
	CreatorProfileID uuid.UUID
}

func (eb ExpenseBill) TableName() string {
	return "group_expense_bills"
}
