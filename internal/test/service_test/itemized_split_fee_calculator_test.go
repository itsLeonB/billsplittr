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

func TestItemizedSplitFeeCalculatorGetMethod(t *testing.T) {
	registry := fee.NewFeeCalculatorRegistry()
	calc := registry[appconstant.ItemizedSplitFee]

	assert.Equal(t, appconstant.ItemizedSplitFee, calc.GetMethod())
}

func TestItemizedSplitFeeCalculatorGetInfo(t *testing.T) {
	registry := fee.NewFeeCalculatorRegistry()
	calc := registry[appconstant.ItemizedSplitFee]

	info := calc.GetInfo()
	assert.Equal(t, appconstant.ItemizedSplitFee, info.Method)
	assert.Equal(t, "Itemized split", info.Display)
	assert.Equal(t, "Split the fee by a fixed rate applied to each expense items", info.Description)
}

func TestItemizedSplitFeeCalculatorValidate(t *testing.T) {
	registry := fee.NewFeeCalculatorRegistry()
	calc := registry[appconstant.ItemizedSplitFee]

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
				Subtotal: decimal.NewFromFloat(100.0),
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
				Subtotal: decimal.NewFromFloat(100.0),
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
				Subtotal: decimal.NewFromFloat(100.0),
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
				Subtotal: decimal.NewFromFloat(100.0),
				Participants: []entity.ExpenseParticipant{
					{ParticipantProfileID: uuid.New()},
				},
			},
			expectError: true,
		},
		{
			name: "zero subtotal",
			fee: entity.OtherFee{
				BaseEntity:     crud.BaseEntity{ID: uuid.New()},
				GroupExpenseID: uuid.New(),
				Amount:         decimal.NewFromFloat(10.0),
			},
			expense: entity.GroupExpense{
				Subtotal: decimal.Zero,
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
				Subtotal:     decimal.NewFromFloat(100.0),
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

func TestItemizedSplitFeeCalculatorSplit(t *testing.T) {
	registry := fee.NewFeeCalculatorRegistry()
	calc := registry[appconstant.ItemizedSplitFee]

	feeID := uuid.New()
	profileID1 := uuid.New()
	profileID2 := uuid.New()

	fee := entity.OtherFee{
		BaseEntity: crud.BaseEntity{ID: feeID},
		Amount:     decimal.NewFromFloat(10.0), // 10% of subtotal
	}

	expense := entity.GroupExpense{
		Subtotal: decimal.NewFromFloat(100.0),
		Participants: []entity.ExpenseParticipant{
			{
				ParticipantProfileID: profileID1,
				ShareAmount:          decimal.NewFromFloat(60.0),
			},
			{
				ParticipantProfileID: profileID2,
				ShareAmount:          decimal.NewFromFloat(40.0),
			},
		},
	}

	result := calc.Split(fee, expense)

	assert.Len(t, result, 2)

	// Rate should be 10/100 = 0.1
	// Participant 1: 60 * 0.1 = 6
	// Participant 2: 40 * 0.1 = 4
	expectedAmounts := map[uuid.UUID]decimal.Decimal{
		profileID1: decimal.NewFromFloat(6.0),
		profileID2: decimal.NewFromFloat(4.0),
	}

	for _, participant := range result {
		assert.Equal(t, feeID, participant.OtherFeeID)
		expectedAmount := expectedAmounts[participant.ProfileID]
		assert.True(t, participant.ShareAmount.Equal(expectedAmount))
	}
}
