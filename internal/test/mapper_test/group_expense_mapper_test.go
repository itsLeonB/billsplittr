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

func TestGroupExpenseRequestToEntity(t *testing.T) {
	creatorID := uuid.New()
	payerID := uuid.New()

	request := dto.NewGroupExpenseRequest{
		CreatorProfileID: creatorID,
		PayerProfileID:   payerID,
		TotalAmount:      decimal.NewFromFloat(100.50),
		Subtotal:         decimal.NewFromFloat(90.00),
		Description:      "Test expense",
		Items: []dto.ExpenseItemData{
			{
				Name:     "Test Item",
				Amount:   decimal.NewFromFloat(45.00),
				Quantity: 2,
			},
		},
		OtherFees: []dto.OtherFeeData{},
	}

	result := mapper.GroupExpenseRequestToEntity(request)

	assert.Equal(t, creatorID, result.CreatorProfileID)
	assert.Equal(t, payerID, result.PayerProfileID)
	assert.True(t, decimal.NewFromFloat(100.50).Equal(result.TotalAmount))
	assert.True(t, decimal.NewFromFloat(90.00).Equal(result.Subtotal))
	assert.Equal(t, "Test expense", result.Description)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, "Test Item", result.Items[0].Name)
	assert.Len(t, result.OtherFees, 0)
}

func TestGroupExpenseToResponse(t *testing.T) {
	id := uuid.New()
	creatorID := uuid.New()
	payerID := uuid.New()
	participantID := uuid.New()
	now := time.Now()

	groupExpense := entity.GroupExpense{
		BaseEntity: crud.BaseEntity{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		PayerProfileID:   payerID,
		TotalAmount:      decimal.NewFromFloat(100.50),
		Subtotal:         decimal.NewFromFloat(90.00),
		Description:      "Test expense",
		CreatorProfileID: creatorID,
		Confirmed:        true,
		Items:            []entity.ExpenseItem{},
		OtherFees:        []entity.OtherFee{},
		Participants: []entity.ExpenseParticipant{
			{
				ParticipantProfileID: participantID,
				ShareAmount:          decimal.NewFromFloat(50.25),
			},
		},
	}

	result := mapper.GroupExpenseToResponse(groupExpense)

	assert.Equal(t, id, result.ID)
	assert.Equal(t, payerID, result.PayerProfileID)
	assert.Equal(t, creatorID, result.CreatorProfileID)
	assert.True(t, decimal.NewFromFloat(100.50).Equal(result.TotalAmount))
	assert.True(t, decimal.NewFromFloat(90.00).Equal(result.Subtotal))
	assert.Equal(t, "Test expense", result.Description)
	assert.True(t, result.Confirmed)
	assert.False(t, result.ParticipantsConfirmed)
	assert.Len(t, result.Participants, 1)
	assert.Equal(t, participantID, result.Participants[0].ProfileID)
	assert.True(t, decimal.NewFromFloat(50.25).Equal(result.Participants[0].ShareAmount))
}
