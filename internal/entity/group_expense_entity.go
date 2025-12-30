package entity

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/go-crud"
	"github.com/shopspring/decimal"
)

type GroupExpense struct {
	crud.BaseEntity
	PayerProfileID uuid.UUID
	TotalAmount    decimal.Decimal
	// Deprecated: use ItemsTotal instead
	Subtotal    decimal.Decimal
	ItemsTotal  decimal.Decimal
	FeesTotal   decimal.Decimal
	Description string
	// Deprecated: refer to Status instead
	Confirmed        bool
	Status           appconstant.ExpenseStatus
	CreatorProfileID uuid.UUID

	// Relationships
	Items        []ExpenseItem        `gorm:"foreignKey:GroupExpenseID"`
	OtherFees    []OtherFee           `gorm:"foreignKey:GroupExpenseID"`
	Participants []ExpenseParticipant `gorm:"foreignKey:GroupExpenseID"`
	Bill         ExpenseBill          `gorm:"foreignKey:GroupExpenseID"`
}

type ExpenseParticipant struct {
	crud.BaseEntity
	GroupExpenseID       uuid.UUID
	ParticipantProfileID uuid.UUID
	ShareAmount          decimal.Decimal
}

func (ep ExpenseParticipant) TableName() string {
	return "group_expense_participants"
}
