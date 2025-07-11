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

func (ge GroupExpense) SimpleName() string {
	return "group expense"
}

type ExpenseItem struct {
	BaseEntity
	GroupExpenseID uuid.UUID
	Name           string
	Amount         decimal.Decimal
	Quantity       int
	Participants   []ItemParticipant `gorm:"foreignKey:ExpenseItemID"`
}

func (ei ExpenseItem) TableName() string {
	return "group_expense_items"
}

func (ei ExpenseItem) SimpleName() string {
	return "expense item"
}

type ItemParticipant struct {
	BaseEntity
	ExpenseItemID uuid.UUID
	ProfileID     uuid.UUID
	Share         decimal.Decimal
}

func (ip ItemParticipant) TableName() string {
	return "group_expense_item_participants"
}

type OtherFee struct {
	BaseEntity
	GroupExpenseID uuid.UUID
	Name           string
	Amount         decimal.Decimal
}

func (of OtherFee) TableName() string {
	return "group_expense_other_fees"
}

type GroupExpenseSpecification struct {
	GroupExpense
	Specification
}
