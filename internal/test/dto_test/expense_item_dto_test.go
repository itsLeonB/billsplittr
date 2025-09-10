package dto_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewExpenseItemRequestFields(t *testing.T) {
	profileID := uuid.New()
	groupExpenseID := uuid.New()

	req := dto.NewExpenseItemRequest{
		ProfileID:      profileID,
		GroupExpenseID: groupExpenseID,
		ExpenseItemData: dto.ExpenseItemData{
			Name:     "Test Item",
			Amount:   decimal.NewFromFloat(25.50),
			Quantity: 2,
		},
	}

	assert.Equal(t, profileID, req.ProfileID)
	assert.Equal(t, groupExpenseID, req.GroupExpenseID)
	assert.Equal(t, "Test Item", req.Name)
	assert.True(t, decimal.NewFromFloat(25.50).Equal(req.Amount))
	assert.Equal(t, 2, req.Quantity)
}

func TestUpdateExpenseItemRequestFields(t *testing.T) {
	profileID := uuid.New()
	itemID := uuid.New()
	groupExpenseID := uuid.New()
	participantID := uuid.New()

	req := dto.UpdateExpenseItemRequest{
		ProfileID:      profileID,
		ID:             itemID,
		GroupExpenseID: groupExpenseID,
		ExpenseItemData: dto.ExpenseItemData{
			Name:     "Updated Item",
			Amount:   decimal.NewFromFloat(30.00),
			Quantity: 3,
		},
		Participants: []dto.ItemParticipantData{
			{
				ProfileID: participantID,
				Share:     decimal.NewFromFloat(15.00),
			},
		},
	}

	assert.Equal(t, profileID, req.ProfileID)
	assert.Equal(t, itemID, req.ID)
	assert.Equal(t, groupExpenseID, req.GroupExpenseID)
	assert.Equal(t, "Updated Item", req.Name)
	assert.Len(t, req.Participants, 1)
	assert.Equal(t, participantID, req.Participants[0].ProfileID)
}

func TestExpenseItemResponseFields(t *testing.T) {
	id := uuid.New()
	groupExpenseID := uuid.New()
	now := time.Now()

	resp := dto.ExpenseItemResponse{
		ID:             id,
		GroupExpenseID: groupExpenseID,
		Name:           "Test Item",
		Amount:         decimal.NewFromFloat(25.50),
		Quantity:       2,
		CreatedAt:      now,
		UpdatedAt:      now,
		Participants:   []dto.ItemParticipantData{},
	}

	assert.Equal(t, id, resp.ID)
	assert.Equal(t, groupExpenseID, resp.GroupExpenseID)
	assert.Equal(t, "Test Item", resp.Name)
	assert.True(t, decimal.NewFromFloat(25.50).Equal(resp.Amount))
	assert.Equal(t, 2, resp.Quantity)
}

func TestItemParticipantDataFields(t *testing.T) {
	profileID := uuid.New()
	share := decimal.NewFromFloat(12.75)

	data := dto.ItemParticipantData{
		ProfileID: profileID,
		Share:     share,
	}

	assert.Equal(t, profileID, data.ProfileID)
	assert.True(t, share.Equal(data.Share))
}

func TestExpenseItemDataFields(t *testing.T) {
	data := dto.ExpenseItemData{
		Name:     "Test Item",
		Amount:   decimal.NewFromFloat(25.50),
		Quantity: 2,
	}

	assert.Equal(t, "Test Item", data.Name)
	assert.True(t, decimal.NewFromFloat(25.50).Equal(data.Amount))
	assert.Equal(t, 2, data.Quantity)
}
