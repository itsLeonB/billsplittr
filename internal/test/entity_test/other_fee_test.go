package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestOtherFee_ProfileIDs(t *testing.T) {
	participantID1 := uuid.New()
	participantID2 := uuid.New()

	fee := entity.OtherFee{
		Name:   "Service Fee",
		Amount: decimal.NewFromFloat(10.0),
		Participants: []entity.FeeParticipant{
			{ProfileID: participantID1},
			{ProfileID: participantID2},
		},
	}

	profileIDs := fee.ProfileIDs()

	assert.Contains(t, profileIDs, participantID1)
	assert.Contains(t, profileIDs, participantID2)
	assert.Len(t, profileIDs, 2)
}
