package helper_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/helper"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGetAmountSumsFromDebtTransactions(t *testing.T) {
	userID := uuid.New()
	friendID := uuid.New()
	otherID := uuid.New()

	tests := []struct {
		name         string
		transactions []entity.DebtTransaction
		wantUser     decimal.Decimal
		wantFriend   decimal.Decimal
	}{
		{
			name:         "Empty transactions",
			transactions: []entity.DebtTransaction{},
			wantUser:     decimal.Zero,
			wantFriend:   decimal.Zero,
		},
		{
			name: "Unrelated transactions",
			transactions: []entity.DebtTransaction{
				{
					LenderProfileID:   otherID,
					BorrowerProfileID: otherID,
					Type:              appconstant.Lend,
					Amount:            decimal.NewFromInt(100),
				},
			},
			wantUser:   decimal.Zero,
			wantFriend: decimal.Zero,
		},
		{
			name: "User lends to friend and gets repaid",
			transactions: []entity.DebtTransaction{
				{
					LenderProfileID:   userID,
					BorrowerProfileID: friendID,
					Type:              appconstant.Lend,
					Amount:            decimal.NewFromInt(100),
				},
				{
					LenderProfileID:   userID,
					BorrowerProfileID: friendID,
					Type:              appconstant.Repay,
					Amount:            decimal.NewFromInt(50),
				},
			},
			wantUser:   decimal.Zero,
			wantFriend: decimal.NewFromInt(50),
		},
		{
			name: "Friend lends to user and gets repaid",
			transactions: []entity.DebtTransaction{
				{
					LenderProfileID:   friendID,
					BorrowerProfileID: userID,
					Type:              appconstant.Lend,
					Amount:            decimal.NewFromInt(70),
				},
				{
					LenderProfileID:   friendID,
					BorrowerProfileID: userID,
					Type:              appconstant.Repay,
					Amount:            decimal.NewFromInt(30),
				},
			},
			wantUser:   decimal.NewFromInt(40),
			wantFriend: decimal.Zero,
		},
		{
			name: "All Lend transactions",
			transactions: []entity.DebtTransaction{
				{
					LenderProfileID:   userID,
					BorrowerProfileID: friendID,
					Type:              appconstant.Lend,
					Amount:            decimal.NewFromInt(20),
				},
				{
					LenderProfileID:   friendID,
					BorrowerProfileID: userID,
					Type:              appconstant.Lend,
					Amount:            decimal.NewFromInt(40),
				},
			},
			wantUser:   decimal.NewFromInt(40), // friend lent 40 to user
			wantFriend: decimal.NewFromInt(20), // user lent 20 to friend
		},
		{
			name: "All Repay transactions",
			transactions: []entity.DebtTransaction{
				{
					LenderProfileID:   userID,
					BorrowerProfileID: friendID,
					Type:              appconstant.Repay,
					Amount:            decimal.NewFromInt(25),
				},
				{
					LenderProfileID:   friendID,
					BorrowerProfileID: userID,
					Type:              appconstant.Repay,
					Amount:            decimal.NewFromInt(35),
				},
			},
			wantUser:   decimal.NewFromInt(-35), // friend repaid 35
			wantFriend: decimal.NewFromInt(-25), // user repaid 25
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, gotFriend := helper.GetAmountSumsFromDebtTransactions(userID, friendID, tt.transactions)

			assert.True(t, gotUser.Equal(tt.wantUser), "UserAmount mismatch: got %v want %v", gotUser, tt.wantUser)
			assert.True(t, gotFriend.Equal(tt.wantFriend), "FriendAmount mismatch: got %v want %v", gotFriend, tt.wantFriend)
		})
	}
}
