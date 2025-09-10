package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestOtherFeeProfileIDs(t *testing.T) {
	profileID1 := uuid.New()
	profileID2 := uuid.New()

	of := entity.OtherFee{
		Participants: []entity.FeeParticipant{
			{ProfileID: profileID1},
			{ProfileID: profileID2},
		},
	}

	profileIDs := of.ProfileIDs()
	assert.Len(t, profileIDs, 2)
	assert.Contains(t, profileIDs, profileID1)
	assert.Contains(t, profileIDs, profileID2)
}

func TestOtherFeeTableName(t *testing.T) {
	of := entity.OtherFee{}
	assert.Equal(t, "group_expense_other_fees", of.TableName())
}

func TestFeeParticipantTableName(t *testing.T) {
	fp := entity.FeeParticipant{}
	assert.Equal(t, "group_expense_other_fee_participants", fp.TableName())
}
