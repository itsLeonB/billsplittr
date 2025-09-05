package entity

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"github.com/shopspring/decimal"
)

type OtherFee struct {
	crud.BaseEntity
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

type FeeParticipant struct {
	crud.BaseEntity
	OtherFeeID  uuid.UUID
	ProfileID   uuid.UUID
	ShareAmount decimal.Decimal
}

func (fp FeeParticipant) TableName() string {
	return "group_expense_other_fee_participants"
}
