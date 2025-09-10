package service_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/service/fee"
	"github.com/itsLeonB/go-crud"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestEqualSplitFeeCalculatorGetMethod(t *testing.T) {
	registry := fee.NewFeeCalculatorRegistry()
	calc := registry[appconstant.EqualSplitFee]

	assert.Equal(t, appconstant.EqualSplitFee, calc.GetMethod())
}

func TestEqualSplitFeeCalculatorGetInfo(t *testing.T) {
	registry := fee.NewFeeCalculatorRegistry()
	calc := registry[appconstant.EqualSplitFee]

	info := calc.GetInfo()
	assert.Equal(t, appconstant.EqualSplitFee, info.Method)
	assert.Equal(t, "Equal split", info.Display)
	assert.Equal(t, "Equally split the fee to all participants", info.Description)
}

func TestEqualSplitFeeCalculatorValidate(t *testing.T) {
	registry := fee.NewFeeCalculatorRegistry()
	calc := registry[appconstant.EqualSplitFee]

	tests := []struct {
		name        string
		fee         entity.OtherFee
		expense     entity.GroupExpense
		expectError bool
	}{
		{
			name: "valid fee and expense",
			fee: entity.OtherFee{
				BaseEntity:     crud.BaseEntity{ID: uuid.New()},
				GroupExpenseID: uuid.New(),
				Amount:         decimal.NewFromFloat(10.0),
			},
			expense: entity.GroupExpense{
				Participants: []entity.ExpenseParticipant{
					{ParticipantProfileID: uuid.New()},
				},
			},
			expectError: false,
		},
		{
			name: "nil fee ID",
			fee: entity.OtherFee{
				BaseEntity:     crud.BaseEntity{ID: uuid.Nil},
				GroupExpenseID: uuid.New(),
				Amount:         decimal.NewFromFloat(10.0),
			},
			expense: entity.GroupExpense{
				Participants: []entity.ExpenseParticipant{
					{ParticipantProfileID: uuid.New()},
				},
			},
			expectError: true,
		},
		{
			name: "nil group expense ID",
			fee: entity.OtherFee{
				BaseEntity:     crud.BaseEntity{ID: uuid.New()},
				GroupExpenseID: uuid.Nil,
				Amount:         decimal.NewFromFloat(10.0),
			},
			expense: entity.GroupExpense{
				Participants: []entity.ExpenseParticipant{
					{ParticipantProfileID: uuid.New()},
				},
			},
			expectError: true,
		},
		{
			name: "zero amount",
			fee: entity.OtherFee{
				BaseEntity:     crud.BaseEntity{ID: uuid.New()},
				GroupExpenseID: uuid.New(),
				Amount:         decimal.Zero,
			},
			expense: entity.GroupExpense{
				Participants: []entity.ExpenseParticipant{
					{ParticipantProfileID: uuid.New()},
				},
			},
			expectError: true,
		},
		{
			name: "no participants",
			fee: entity.OtherFee{
				BaseEntity:     crud.BaseEntity{ID: uuid.New()},
				GroupExpenseID: uuid.New(),
				Amount:         decimal.NewFromFloat(10.0),
			},
			expense: entity.GroupExpense{
				Participants: []entity.ExpenseParticipant{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := calc.Validate(tt.fee, tt.expense)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEqualSplitFeeCalculatorSplit(t *testing.T) {
	registry := fee.NewFeeCalculatorRegistry()
	calc := registry[appconstant.EqualSplitFee]

	feeID := uuid.New()
	profileID1 := uuid.New()
	profileID2 := uuid.New()

	fee := entity.OtherFee{
		BaseEntity: crud.BaseEntity{ID: feeID},
		Amount:     decimal.NewFromFloat(20.0),
	}

	expense := entity.GroupExpense{
		Participants: []entity.ExpenseParticipant{
			{ParticipantProfileID: profileID1},
			{ParticipantProfileID: profileID2},
		},
	}

	result := calc.Split(fee, expense)

	assert.Len(t, result, 2)

	expectedAmount := decimal.NewFromFloat(10.0)
	for _, participant := range result {
		assert.Equal(t, feeID, participant.OtherFeeID)
		assert.True(t, participant.ShareAmount.Equal(expectedAmount))
		assert.True(t, participant.ProfileID == profileID1 || participant.ProfileID == profileID2)
	}
}
