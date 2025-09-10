package mapper_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/mapper"
	"github.com/itsLeonB/go-crud"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestExpenseItemToResponse(t *testing.T) {
	id := uuid.New()
	groupExpenseID := uuid.New()
	participantID := uuid.New()
	now := time.Now()

	item := entity.ExpenseItem{
		BaseEntity: crud.BaseEntity{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		GroupExpenseID: groupExpenseID,
		Name:           "Test Item",
		Amount:         decimal.NewFromFloat(25.50),
		Quantity:       2,
		Participants: []entity.ItemParticipant{
			{
				ProfileID: participantID,
				Share:     decimal.NewFromFloat(12.75),
			},
		},
	}

	result := mapper.ExpenseItemToResponse(item)

	assert.Equal(t, id, result.ID)
	assert.Equal(t, groupExpenseID, result.GroupExpenseID)
	assert.Equal(t, "Test Item", result.Name)
	assert.True(t, decimal.NewFromFloat(25.50).Equal(result.Amount))
	assert.Equal(t, 2, result.Quantity)
	assert.Len(t, result.Participants, 1)
	assert.Equal(t, participantID, result.Participants[0].ProfileID)
	assert.True(t, decimal.NewFromFloat(12.75).Equal(result.Participants[0].Share))
}

func TestExpenseItemRequestToEntity(t *testing.T) {
	profileID := uuid.New()
	groupExpenseID := uuid.New()

	request := dto.NewExpenseItemRequest{
		ProfileID:      profileID,
		GroupExpenseID: groupExpenseID,
		ExpenseItemData: dto.ExpenseItemData{
			Name:     "Test Item",
			Amount:   decimal.NewFromFloat(25.50),
			Quantity: 2,
		},
	}

	result := mapper.ExpenseItemRequestToEntity(request)

	assert.Equal(t, groupExpenseID, result.GroupExpenseID)
	assert.Equal(t, "Test Item", result.Name)
	assert.True(t, decimal.NewFromFloat(25.50).Equal(result.Amount))
	assert.Equal(t, 2, result.Quantity)
}

func TestPatchExpenseItemWithRequest(t *testing.T) {
	id := uuid.New()
	profileID := uuid.New()
	groupExpenseID := uuid.New()

	expenseItem := entity.ExpenseItem{
		BaseEntity: crud.BaseEntity{ID: id},
		Name:       "Old Item",
		Amount:     decimal.NewFromFloat(10.00),
		Quantity:   1,
	}

	request := dto.UpdateExpenseItemRequest{
		ProfileID:      profileID,
		ID:             id,
		GroupExpenseID: groupExpenseID,
		ExpenseItemData: dto.ExpenseItemData{
			Name:     "Updated Item",
			Amount:   decimal.NewFromFloat(30.00),
			Quantity: 3,
		},
	}

	result := mapper.PatchExpenseItemWithRequest(expenseItem, request)

	assert.Equal(t, id, result.ID)
	assert.Equal(t, "Updated Item", result.Name)
	assert.True(t, decimal.NewFromFloat(30.00).Equal(result.Amount))
	assert.Equal(t, 3, result.Quantity)
}

func TestItemParticipantRequestToEntity(t *testing.T) {
	profileID := uuid.New()

	data := dto.ItemParticipantData{
		ProfileID: profileID,
		Share:     decimal.NewFromFloat(15.00),
	}

	result := mapper.ItemParticipantRequestToEntity(data)

	assert.Equal(t, profileID, result.ProfileID)
	assert.True(t, decimal.NewFromFloat(15.00).Equal(result.Share))
}
