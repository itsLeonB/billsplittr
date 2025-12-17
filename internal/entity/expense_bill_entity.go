package entity

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/go-crud"
)

type ExpenseBill struct {
	crud.BaseEntity
	GroupExpenseID   uuid.NullUUID
	PayerProfileID   uuid.UUID
	CreatorProfileID uuid.UUID
	ImageName        string
	Status           appconstant.BillStatus
}

func (eb ExpenseBill) TableName() string {
	return "group_expense_bills"
}
