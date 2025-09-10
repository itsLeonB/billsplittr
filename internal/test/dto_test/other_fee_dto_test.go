package dto_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewOtherFeeRequestFields(t *testing.T) {
	profileID := uuid.New()
	groupExpenseID := uuid.New()

	req := dto.NewOtherFeeRequest{
		ProfileID:      profileID,
		GroupExpenseID: groupExpenseID,
		OtherFeeData: dto.OtherFeeData{
			Name:              "Service Fee",
			Amount:            decimal.NewFromFloat(10.50),
			CalculationMethod: appconstant.EqualSplitFee,
		},
	}

	assert.Equal(t, profileID, req.ProfileID)
	assert.Equal(t, groupExpenseID, req.GroupExpenseID)
	assert.Equal(t, "Service Fee", req.Name)
	assert.True(t, decimal.NewFromFloat(10.50).Equal(req.Amount))
	assert.Equal(t, appconstant.EqualSplitFee, req.CalculationMethod)
}

func TestUpdateOtherFeeRequestFields(t *testing.T) {
	profileID := uuid.New()
	feeID := uuid.New()
	groupExpenseID := uuid.New()

	req := dto.UpdateOtherFeeRequest{
		ProfileID:      profileID,
		ID:             feeID,
		GroupExpenseID: groupExpenseID,
		OtherFeeData: dto.OtherFeeData{
			Name:              "Updated Fee",
			Amount:            decimal.NewFromFloat(15.00),
			CalculationMethod: appconstant.ItemizedSplitFee,
		},
	}

	assert.Equal(t, profileID, req.ProfileID)
	assert.Equal(t, feeID, req.ID)
	assert.Equal(t, groupExpenseID, req.GroupExpenseID)
	assert.Equal(t, "Updated Fee", req.Name)
	assert.Equal(t, appconstant.ItemizedSplitFee, req.CalculationMethod)
}

func TestOtherFeeResponseFields(t *testing.T) {
	id := uuid.New()
	groupExpenseID := uuid.New()
	now := time.Now()

	resp := dto.OtherFeeResponse{
		ID:                id,
		GroupExpenseID:    groupExpenseID,
		Name:              "Service Fee",
		Amount:            decimal.NewFromFloat(10.50),
		CalculationMethod: appconstant.EqualSplitFee,
		CreatedAt:         now,
		UpdatedAt:         now,
		Participants:      []dto.FeeParticipantResponse{},
	}

	assert.Equal(t, id, resp.ID)
	assert.Equal(t, groupExpenseID, resp.GroupExpenseID)
	assert.Equal(t, "Service Fee", resp.Name)
	assert.True(t, decimal.NewFromFloat(10.50).Equal(resp.Amount))
	assert.Equal(t, appconstant.EqualSplitFee, resp.CalculationMethod)
}

func TestFeeParticipantResponseFields(t *testing.T) {
	profileID := uuid.New()
	shareAmount := decimal.NewFromFloat(5.25)

	resp := dto.FeeParticipantResponse{
		ProfileID:   profileID,
		ShareAmount: shareAmount,
	}

	assert.Equal(t, profileID, resp.ProfileID)
	assert.True(t, shareAmount.Equal(resp.ShareAmount))
}

func TestFeeCalculationMethodInfoFields(t *testing.T) {
	info := dto.FeeCalculationMethodInfo{
		Method:      appconstant.EqualSplitFee,
		Display:     "Fixed Amount",
		Description: "A fixed fee amount",
	}

	assert.Equal(t, appconstant.EqualSplitFee, info.Method)
	assert.Equal(t, "Fixed Amount", info.Display)
	assert.Equal(t, "A fixed fee amount", info.Description)
}

func TestOtherFeeDataFields(t *testing.T) {
	data := dto.OtherFeeData{
		Name:              "Service Fee",
		Amount:            decimal.NewFromFloat(10.50),
		CalculationMethod: appconstant.EqualSplitFee,
	}

	assert.Equal(t, "Service Fee", data.Name)
	assert.True(t, decimal.NewFromFloat(10.50).Equal(data.Amount))
	assert.Equal(t, appconstant.EqualSplitFee, data.CalculationMethod)
}
