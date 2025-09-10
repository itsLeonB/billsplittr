package mapper_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/mapper"
	"github.com/itsLeonB/go-crud"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestOtherFeeToResponse(t *testing.T) {
	id := uuid.New()
	groupExpenseID := uuid.New()
	participantID := uuid.New()
	now := time.Now()

	fee := entity.OtherFee{
		BaseEntity: crud.BaseEntity{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		GroupExpenseID:    groupExpenseID,
		Name:              "Service Fee",
		Amount:            decimal.NewFromFloat(10.50),
		CalculationMethod: appconstant.EqualSplitFee,
		Participants: []entity.FeeParticipant{
			{
				ProfileID:   participantID,
				ShareAmount: decimal.NewFromFloat(5.25),
			},
		},
	}

	result := mapper.OtherFeeToResponse(fee)

	assert.Equal(t, id, result.ID)
	assert.Equal(t, groupExpenseID, result.GroupExpenseID)
	assert.Equal(t, "Service Fee", result.Name)
	assert.True(t, decimal.NewFromFloat(10.50).Equal(result.Amount))
	assert.Equal(t, appconstant.EqualSplitFee, result.CalculationMethod)
	assert.Len(t, result.Participants, 1)
	assert.Equal(t, participantID, result.Participants[0].ProfileID)
	assert.True(t, decimal.NewFromFloat(5.25).Equal(result.Participants[0].ShareAmount))
}

func TestOtherFeeRequestToEntity(t *testing.T) {
	request := dto.OtherFeeData{
		Name:              "Service Fee",
		Amount:            decimal.NewFromFloat(10.50),
		CalculationMethod: appconstant.EqualSplitFee,
	}

	result := mapper.OtherFeeRequestToEntity(request)

	assert.Equal(t, "Service Fee", result.Name)
	assert.True(t, decimal.NewFromFloat(10.50).Equal(result.Amount))
	assert.Equal(t, appconstant.EqualSplitFee, result.CalculationMethod)
}

func TestPatchOtherFeeWithRequest(t *testing.T) {
	id := uuid.New()
	profileID := uuid.New()
	groupExpenseID := uuid.New()

	otherFee := entity.OtherFee{
		BaseEntity:        crud.BaseEntity{ID: id},
		Name:              "Old Fee",
		Amount:            decimal.NewFromFloat(5.00),
		CalculationMethod: appconstant.EqualSplitFee,
	}

	request := dto.UpdateOtherFeeRequest{
		ProfileID:      profileID,
		ID:             id,
		GroupExpenseID: groupExpenseID,
		OtherFeeData: dto.OtherFeeData{
			Name:              "Updated Fee",
			Amount:            decimal.NewFromFloat(15.00),
			CalculationMethod: appconstant.ItemizedSplitFee,
		},
	}

	result := mapper.PatchOtherFeeWithRequest(otherFee, request)

	assert.Equal(t, id, result.ID)
	assert.Equal(t, "Updated Fee", result.Name)
	assert.True(t, decimal.NewFromFloat(15.00).Equal(result.Amount))
	assert.Equal(t, appconstant.ItemizedSplitFee, result.CalculationMethod)
}
