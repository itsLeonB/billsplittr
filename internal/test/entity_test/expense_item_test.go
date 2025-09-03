package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestExpenseItem_ProfileIDs(t *testing.T) {
	participantID1 := uuid.New()
	participantID2 := uuid.New()

	item := entity.ExpenseItem{
		Name:     "Test Item",
		Amount:   decimal.NewFromFloat(50.0),
		Quantity: 2,
		Participants: []entity.ItemParticipant{
			{ProfileID: participantID1},
			{ProfileID: participantID2},
		},
	}

	profileIDs := item.ProfileIDs()

	assert.Contains(t, profileIDs, participantID1)
	assert.Contains(t, profileIDs, participantID2)
	assert.Len(t, profileIDs, 2)
}

func TestExpenseItem_TotalAmount(t *testing.T) {
	item := entity.ExpenseItem{
		Amount:   decimal.NewFromFloat(25.0),
		Quantity: 3,
	}

	total := item.TotalAmount()
	expected := decimal.NewFromFloat(75.0)

	assert.True(t, total.Equal(expected))
}
