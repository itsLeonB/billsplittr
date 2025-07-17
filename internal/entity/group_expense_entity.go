package entity

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/shopspring/decimal"
)

type GroupExpense struct {
	BaseEntity
	PayerProfileID        uuid.UUID
	TotalAmount           decimal.Decimal
	Description           string
	Items                 []ExpenseItem `gorm:"foreignKey:GroupExpenseID"`
	OtherFees             []OtherFee    `gorm:"foreignKey:GroupExpenseID"`
	Confirmed             bool
	ParticipantsConfirmed bool
	CreatorProfileID      uuid.UUID
	PayerProfile          UserProfile          `gorm:"foreignKey:PayerProfileID"`
	CreatorProfile        UserProfile          `gorm:"foreignKey:CreatorProfileID"`
	Participants          []ExpenseParticipant `gorm:"foreignKey:GroupExpenseID"`
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
	Profile       UserProfile `gorm:"foreignKey:ProfileID"`
}

func (ip ItemParticipant) TableName() string {
	return "group_expense_item_participants"
}

type OtherFee struct {
	BaseEntity
	GroupExpenseID    uuid.UUID
	Name              string
	Amount            decimal.Decimal
	CalculationMethod appconstant.FeeCalculationMethod
	Rate              decimal.NullDecimal
	Participants      []FeeParticipant `gorm:"foreignKey:OtherFeeID"`
}

func (of OtherFee) TableName() string {
	return "group_expense_other_fees"
}

type GroupExpenseSpecification struct {
	GroupExpense
	Specification
}

type FeeParticipant struct {
	BaseEntity
	OtherFeeID  uuid.UUID
	ProfileID   uuid.UUID
	ShareAmount decimal.Decimal
	Profile     UserProfile `gorm:"foreignKey:ProfileID"`
}

func (fp FeeParticipant) TableName() string {
	return "group_expense_other_fee_participants"
}
