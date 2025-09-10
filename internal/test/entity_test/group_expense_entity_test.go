package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestGroupExpenseProfileIDs(t *testing.T) {
	creatorID := uuid.New()
	payerID := uuid.New()
	participantID := uuid.New()
	itemParticipantID := uuid.New()
	feeParticipantID := uuid.New()

	ge := entity.GroupExpense{
		CreatorProfileID: creatorID,
		PayerProfileID:   payerID,
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

	profileIDs := ge.ProfileIDs()

	assert.Contains(t, profileIDs, creatorID)
	assert.Contains(t, profileIDs, payerID)
	assert.Contains(t, profileIDs, participantID)
	assert.Contains(t, profileIDs, itemParticipantID)
	assert.Contains(t, profileIDs, feeParticipantID)
}

func TestExpenseParticipantTableName(t *testing.T) {
	ep := entity.ExpenseParticipant{}
	assert.Equal(t, "group_expense_participants", ep.TableName())
}
