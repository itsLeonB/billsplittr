package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type GroupExpense struct {
	BaseEntity
	PayerProfileID        uuid.UUID
	TotalAmount           decimal.Decimal
	Description           string
	Items                 []ExpenseItem
	OtherFees             []OtherFee
	Confirmed             bool
	ParticipantsConfirmed bool
	CreatorProfileID      uuid.UUID
	PayerProfile          UserProfile `gorm:"foreignKey:PayerProfileID"`
	CreatorProfile        UserProfile `gorm:"foreignKey:CreatorProfileID"`
}

type ExpenseItem struct {
	BaseEntity
	GroupExpenseID uuid.UUID
	Name           string
	Amount         decimal.Decimal
	Quantity       int
}

func (ei *ExpenseItem) TableName() string {
	return "group_expense_items"
}

type OtherFee struct {
	BaseEntity
	GroupExpenseID uuid.UUID
	Name           string
	Amount         decimal.Decimal
}

func (of *OtherFee) TableName() string {
	return "group_expense_other_fees"
}

type GroupExpenseSpecification struct {
	GroupExpense
	Specification
}
