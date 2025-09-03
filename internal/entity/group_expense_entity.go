package entity

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/go-crud"
	"github.com/shopspring/decimal"
)

type GroupExpense struct {
	crud.BaseEntity
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
