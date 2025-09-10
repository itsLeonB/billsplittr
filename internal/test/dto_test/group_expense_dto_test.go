package dto_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewGroupExpenseRequestFields(t *testing.T) {
	creatorID := uuid.New()
	payerID := uuid.New()
	totalAmount := decimal.NewFromFloat(100.50)
	subtotal := decimal.NewFromFloat(90.00)

	req := dto.NewGroupExpenseRequest{
		CreatorProfileID: creatorID,
		PayerProfileID:   payerID,
		TotalAmount:      totalAmount,
		Subtotal:         subtotal,
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

	assert.Equal(t, creatorID, req.CreatorProfileID)
	assert.Equal(t, payerID, req.PayerProfileID)
	assert.True(t, totalAmount.Equal(req.TotalAmount))
	assert.True(t, subtotal.Equal(req.Subtotal))
	assert.Equal(t, "Test expense", req.Description)
	assert.Len(t, req.Items, 1)
	assert.Equal(t, "Test Item", req.Items[0].Name)
}

func TestGroupExpenseResponseFields(t *testing.T) {
	id := uuid.New()
	creatorID := uuid.New()
	payerID := uuid.New()
	now := time.Now()

	resp := dto.GroupExpenseResponse{
		ID:                    id,
		PayerProfileID:        payerID,
		TotalAmount:           decimal.NewFromFloat(100.50),
		Subtotal:              decimal.NewFromFloat(90.00),
		Description:           "Test expense",
		CreatorProfileID:      creatorID,
		Confirmed:             true,
		ParticipantsConfirmed: false,
		CreatedAt:             now,
		UpdatedAt:             now,
		Items:                 []dto.ExpenseItemResponse{},
		OtherFees:             []dto.OtherFeeResponse{},
		Participants:          []dto.ExpenseParticipantResponse{},
	}

	assert.Equal(t, id, resp.ID)
	assert.Equal(t, payerID, resp.PayerProfileID)
	assert.Equal(t, creatorID, resp.CreatorProfileID)
	assert.True(t, resp.Confirmed)
	assert.False(t, resp.ParticipantsConfirmed)
}

func TestExpenseParticipantResponseFields(t *testing.T) {
	profileID := uuid.New()
	shareAmount := decimal.NewFromFloat(25.50)

	resp := dto.ExpenseParticipantResponse{
		ProfileID:   profileID,
		ShareAmount: shareAmount,
	}

	assert.Equal(t, profileID, resp.ProfileID)
	assert.True(t, shareAmount.Equal(resp.ShareAmount))
}
