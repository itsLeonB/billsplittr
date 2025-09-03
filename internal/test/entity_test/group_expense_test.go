package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGroupExpense_ProfileIDs(t *testing.T) {
	creatorID := uuid.New()
	payerID := uuid.New()
	participantID := uuid.New()
	itemParticipantID := uuid.New()
	feeParticipantID := uuid.New()

	groupExpense := entity.GroupExpense{
		CreatorProfileID: creatorID,
		PayerProfileID:   payerID,
		TotalAmount:      decimal.NewFromFloat(100.0),
		Subtotal:         decimal.NewFromFloat(90.0),
		Description:      "Test expense",
		Items: []entity.ExpenseItem{
			{
				Participants: []entity.ItemParticipant{
					{ProfileID: itemParticipantID},
				},
			},
		},
		OtherFees: []entity.OtherFee{
			{
				Participants: []entity.FeeParticipant{
					{ProfileID: feeParticipantID},
				},
			},
		},
		Participants: []entity.ExpenseParticipant{
			{ParticipantProfileID: participantID},
		},
	}

	profileIDs := groupExpense.ProfileIDs()

	assert.Contains(t, profileIDs, creatorID)
	assert.Contains(t, profileIDs, payerID)
	assert.Contains(t, profileIDs, participantID)
	assert.Contains(t, profileIDs, itemParticipantID)
	assert.Contains(t, profileIDs, feeParticipantID)
}
