package entity

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/shopspring/decimal"
)

type GroupExpense struct {
	BaseEntity
	PayerProfileID        uuid.UUID
	TotalAmount           decimal.Decimal
	Subtotal              decimal.Decimal
	Description           string
	Items                 []ExpenseItem `gorm:"foreignKey:GroupExpenseID"`
	OtherFees             []OtherFee    `gorm:"foreignKey:GroupExpenseID"`
	Confirmed             bool
	ParticipantsConfirmed bool
	CreatorProfileID      uuid.UUID
	Participants          []ExpenseParticipant `gorm:"foreignKey:GroupExpenseID"`
}

func (ge GroupExpense) SimpleName() string {
	return "group expense"
}

func (ge GroupExpense) ProfileIDs() []uuid.UUID {
	profileIDs := make([]uuid.UUID, 0)
	profileIDs = append(profileIDs, ge.CreatorProfileID)
	profileIDs = append(profileIDs, ge.PayerProfileID)
	for _, item := range ge.Items {
		profileIDs = append(profileIDs, item.ProfileIDs()...)
	}
	for _, fee := range ge.OtherFees {
		profileIDs = append(profileIDs, fee.ProfileIDs()...)
	}
	for _, participant := range ge.Participants {
		profileIDs = append(profileIDs, participant.ParticipantProfileID)
	}

	return profileIDs
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

func (ei ExpenseItem) TotalAmount() decimal.Decimal {
	return ei.Amount.Mul(decimal.NewFromInt(int64(ei.Quantity)))
}

func (ei ExpenseItem) ProfileIDs() []uuid.UUID {
	return ezutil.MapSlice(ei.Participants, func(ip ItemParticipant) uuid.UUID { return ip.ProfileID })
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
	GroupExpenseID    uuid.UUID
	Name              string
	Amount            decimal.Decimal
	CalculationMethod appconstant.FeeCalculationMethod
	Participants      []FeeParticipant `gorm:"foreignKey:OtherFeeID"`
}

func (of OtherFee) ProfileIDs() []uuid.UUID {
	return ezutil.MapSlice(of.Participants, func(fp FeeParticipant) uuid.UUID { return fp.ProfileID })
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
}

func (fp FeeParticipant) TableName() string {
	return "group_expense_other_fee_participants"
}

type ExpenseBill struct {
	BaseEntity
	PayerProfileID   uuid.UUID
	ImageName        string
	GroupExpenseID   uuid.NullUUID
	CreatorProfileID uuid.UUID
}

func (eb ExpenseBill) TableName() string {
	return "group_expense_bills"
}
