package debt_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/service/debt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestAnonymousDebtCalculatorStrategies(t *testing.T) {
	calculators := debt.NewAnonymousDebtCalculatorStrategies()
	userProfileID := uuid.New()
	friendProfileID := uuid.New()

	tests := []struct {
		name             string
		request          dto.NewDebtTransactionRequest
		existingEntities []entity.DebtTransaction
		expectedEntity   entity.DebtTransaction
		expectedError    bool
	}{
		{
			name: "Borrowing action - simple success",
			request: dto.NewDebtTransactionRequest{
				UserProfileID:    userProfileID,
				FriendProfileID:  friendProfileID,
				Action:           appconstant.BorrowAction,
				Amount:           decimal.NewFromInt(100),
				TransferMethodID: uuid.New(),
				Description:      "Borrowing test",
			},
			expectedEntity: entity.DebtTransaction{
				Type: appconstant.Lend,
			},
			expectedError: false,
		},
		{
			name: "Lending action - simple success",
			request: dto.NewDebtTransactionRequest{
				UserProfileID:    userProfileID,
				FriendProfileID:  friendProfileID,
				Action:           appconstant.LendAction,
				Amount:           decimal.NewFromInt(200),
				TransferMethodID: uuid.New(),
				Description:      "Lending test",
			},
			expectedEntity: entity.DebtTransaction{
				Type: appconstant.Lend,
			},
			expectedError: false,
		},
		{
			name: "Receiving action - equal to limit",
			request: dto.NewDebtTransactionRequest{
				UserProfileID:    userProfileID,
				FriendProfileID:  friendProfileID,
				Action:           appconstant.ReceiveAction,
				Amount:           decimal.NewFromInt(100),
				TransferMethodID: uuid.New(),
				Description:      "Receive equal with limit test",
			},
			existingEntities: []entity.DebtTransaction{
				{
					LenderProfileID:   userProfileID,
					BorrowerProfileID: friendProfileID,
					Type:              appconstant.Lend,
					Amount:            decimal.NewFromInt(100),
				},
			},
			expectedEntity: entity.DebtTransaction{
				Type: appconstant.Repay,
			},
			expectedError: false,
		},
		{
			name: "Receiving action - fail over limit",
			request: dto.NewDebtTransactionRequest{
				UserProfileID:    userProfileID,
				FriendProfileID:  friendProfileID,
				Action:           appconstant.ReceiveAction,
				Amount:           decimal.NewFromInt(500),
				TransferMethodID: uuid.New(),
				Description:      "Receive over limit test",
			},
			existingEntities: []entity.DebtTransaction{
				{
					LenderProfileID:   userProfileID,
					BorrowerProfileID: friendProfileID,
					Type:              appconstant.Lend,
					Amount:            decimal.NewFromInt(100),
				},
			},
			expectedEntity: entity.DebtTransaction{
				Type: appconstant.Repay,
			},
			expectedError: true,
		},
		{
			name: "Receiving action - success under limit",
			request: dto.NewDebtTransactionRequest{
				UserProfileID:    userProfileID,
				FriendProfileID:  friendProfileID,
				Action:           appconstant.ReceiveAction,
				Amount:           decimal.NewFromInt(50),
				TransferMethodID: uuid.New(),
				Description:      "Receive under limit test",
			},
			existingEntities: []entity.DebtTransaction{
				{
					LenderProfileID:   userProfileID,
					BorrowerProfileID: friendProfileID,
					Type:              appconstant.Lend,
					Amount:            decimal.NewFromInt(100),
				},
			},
			expectedEntity: entity.DebtTransaction{
				Type: appconstant.Repay,
			},
			expectedError: false,
		},
		{
			name: "Returning action - equal to limit",
			request: dto.NewDebtTransactionRequest{
				UserProfileID:    userProfileID,
				FriendProfileID:  friendProfileID,
				Action:           appconstant.ReturnAction,
				Amount:           decimal.NewFromInt(100),
				TransferMethodID: uuid.New(),
				Description:      "Return equal with limit test",
			},
			existingEntities: []entity.DebtTransaction{
				{
					LenderProfileID:   friendProfileID,
					BorrowerProfileID: userProfileID,
					Type:              appconstant.Lend,
					Amount:            decimal.NewFromInt(100),
				},
			},
			expectedEntity: entity.DebtTransaction{
				Type: appconstant.Repay,
			},
			expectedError: false,
		},
		{
			name: "Returning action - fail over limit",
			request: dto.NewDebtTransactionRequest{
				UserProfileID:    userProfileID,
				FriendProfileID:  friendProfileID,
				Action:           appconstant.ReturnAction,
				Amount:           decimal.NewFromInt(500),
				TransferMethodID: uuid.New(),
				Description:      "Return test",
			},
			existingEntities: []entity.DebtTransaction{
				{
					LenderProfileID:   friendProfileID,
					BorrowerProfileID: userProfileID,
					Type:              appconstant.Lend,
					Amount:            decimal.NewFromInt(100),
				},
			},
			expectedEntity: entity.DebtTransaction{
				Type: appconstant.Repay,
			},
			expectedError: true,
		},
		{
			name: "Returning action - success under limit",
			request: dto.NewDebtTransactionRequest{
				UserProfileID:    userProfileID,
				FriendProfileID:  friendProfileID,
				Action:           appconstant.ReturnAction,
				Amount:           decimal.NewFromInt(50),
				TransferMethodID: uuid.New(),
				Description:      "Return test",
			},
			existingEntities: []entity.DebtTransaction{
				{
					LenderProfileID:   friendProfileID,
					BorrowerProfileID: userProfileID,
					Type:              appconstant.Lend,
					Amount:            decimal.NewFromInt(100),
				},
			},
			expectedEntity: entity.DebtTransaction{
				Type: appconstant.Repay,
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculator, ok := calculators[tt.request.Action]
			assert.True(t, ok, "Calculator for action should exist")

			// MapRequestToEntity
			entityResult := calculator.MapRequestToEntity(tt.request)
			assert.Equal(t, tt.expectedEntity.Type, entityResult.Type, "Entity Type should match")
			assert.True(t, tt.request.Amount.Compare(entityResult.Amount) == 0, "Amount should match")
			assert.Equal(t, tt.request.Description, entityResult.Description, "Description should match")

			// Validate
			err := calculator.Validate(entityResult, tt.existingEntities)
			if tt.expectedError {
				assert.Error(t, err, "Expected validation error but got none")
			} else {
				assert.NoError(t, err, "Did not expect validation error but got one")
			}

			// MapEntityToResponse
			entityResult.ID = uuid.New()
			entityResult.TransferMethod = entity.TransferMethod{
				Display: "TestMethod",
			}
			entityResult.CreatedAt = time.Now()
			entityResult.UpdatedAt = time.Now()

			response := calculator.MapEntityToResponse(entityResult)
			assert.True(t, entityResult.Amount.Compare(response.Amount) == 0, "Response Amount should match Entity")
			assert.Equal(t, entityResult.Type, response.Type, "Response Type should match Entity")
			assert.Equal(t, entityResult.Description, response.Description, "Response Description should match Entity")
			assert.Equal(t, "TestMethod", response.TransferMethod, "Response TransferMethod should match Entity")
		})
	}
}
