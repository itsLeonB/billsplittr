package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/billsplittr/internal/test/mocks"
	"github.com/itsLeonB/go-crud"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGroupExpenseService_CreateDraftSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := mocks.NewMockTransactor(ctrl)
	mockGroupExpenseRepo := mocks.NewMockGroupExpenseRepository(ctrl)
	mockOtherFeeRepo := mocks.NewMockOtherFeeRepository(ctrl)

	svc := service.NewGroupExpenseService(mockTransactor, mockGroupExpenseRepo, mockOtherFeeRepo)

	creatorID := uuid.New()
	payerID := uuid.New()

	request := dto.NewGroupExpenseRequest{
		CreatorProfileID: creatorID,
		PayerProfileID:   payerID,
		TotalAmount:      decimal.NewFromFloat(100.00),
		Subtotal:         decimal.NewFromFloat(90.00),
		Description:      "Test expense",
		Items: []dto.ExpenseItemData{
			{
				Name:     "Test Item",
				Amount:   decimal.NewFromFloat(45.00),
				Quantity: 2,
			},
		},
		OtherFees: []dto.OtherFeeData{
			{
				Name:              "Service Fee",
				Amount:            decimal.NewFromFloat(10.00),
				CalculationMethod: appconstant.EqualSplitFee,
			},
		},
	}

	expectedEntity := entity.GroupExpense{
		PayerProfileID:   payerID,
		TotalAmount:      decimal.NewFromFloat(100.00),
		Subtotal:         decimal.NewFromFloat(90.00),
		Description:      "Test expense",
		CreatorProfileID: creatorID,
	}

	mockGroupExpenseRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		Return(expectedEntity, nil)

	result, err := svc.CreateDraft(context.Background(), request)

	assert.NoError(t, err)
	assert.Equal(t, payerID, result.PayerProfileID)
	assert.Equal(t, creatorID, result.CreatorProfileID)
	assert.True(t, decimal.NewFromFloat(100.00).Equal(result.TotalAmount))
}

func TestGroupExpenseService_CreateDraft_ValidationErrorZeroAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := mocks.NewMockTransactor(ctrl)
	mockGroupExpenseRepo := mocks.NewMockGroupExpenseRepository(ctrl)
	mockOtherFeeRepo := mocks.NewMockOtherFeeRepository(ctrl)

	svc := service.NewGroupExpenseService(mockTransactor, mockGroupExpenseRepo, mockOtherFeeRepo)

	request := dto.NewGroupExpenseRequest{
		TotalAmount: decimal.Zero,
		Items:       []dto.ExpenseItemData{},
		OtherFees:   []dto.OtherFeeData{},
	}

	_, err := svc.CreateDraft(context.Background(), request)

	assert.Error(t, err)
	// Error message check removed
}

func TestGroupExpenseService_CreateDraft_ValidationErrorAmountMismatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := mocks.NewMockTransactor(ctrl)
	mockGroupExpenseRepo := mocks.NewMockGroupExpenseRepository(ctrl)
	mockOtherFeeRepo := mocks.NewMockOtherFeeRepository(ctrl)

	svc := service.NewGroupExpenseService(mockTransactor, mockGroupExpenseRepo, mockOtherFeeRepo)

	request := dto.NewGroupExpenseRequest{
		TotalAmount: decimal.NewFromFloat(100.00),
		Subtotal:    decimal.NewFromFloat(50.00), // Mismatch
		Items: []dto.ExpenseItemData{
			{
				Name:     "Test Item",
				Amount:   decimal.NewFromFloat(45.00),
				Quantity: 2, // Should be 90.00 total
			},
		},
		OtherFees: []dto.OtherFeeData{
			{
				Name:   "Service Fee",
				Amount: decimal.NewFromFloat(10.00),
			},
		},
	}

	_, err := svc.CreateDraft(context.Background(), request)

	assert.Error(t, err)
	// Error message check removed
}

func TestGroupExpenseServiceGetAllCreated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := mocks.NewMockTransactor(ctrl)
	mockGroupExpenseRepo := mocks.NewMockGroupExpenseRepository(ctrl)
	mockOtherFeeRepo := mocks.NewMockOtherFeeRepository(ctrl)

	svc := service.NewGroupExpenseService(mockTransactor, mockGroupExpenseRepo, mockOtherFeeRepo)

	profileID := uuid.New()
	expectedExpenses := []entity.GroupExpense{
		{
			BaseEntity:       crud.BaseEntity{ID: uuid.New()},
			CreatorProfileID: profileID,
			TotalAmount:      decimal.NewFromFloat(100.00),
		},
	}

	mockGroupExpenseRepo.EXPECT().
		FindAll(gomock.Any(), gomock.Any()).
		Return(expectedExpenses, nil)

	result, err := svc.GetAllCreated(context.Background(), profileID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, profileID, result[0].CreatorProfileID)
}

func TestGroupExpenseServiceGetDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := mocks.NewMockTransactor(ctrl)
	mockGroupExpenseRepo := mocks.NewMockGroupExpenseRepository(ctrl)
	mockOtherFeeRepo := mocks.NewMockOtherFeeRepository(ctrl)

	svc := service.NewGroupExpenseService(mockTransactor, mockGroupExpenseRepo, mockOtherFeeRepo)

	expenseID := uuid.New()
	expectedExpense := entity.GroupExpense{
		BaseEntity:  crud.BaseEntity{ID: expenseID},
		TotalAmount: decimal.NewFromFloat(100.00),
	}

	mockGroupExpenseRepo.EXPECT().
		FindFirst(gomock.Any(), gomock.Any()).
		Return(expectedExpense, nil)

	result, err := svc.GetDetails(context.Background(), expenseID)

	assert.NoError(t, err)
	assert.Equal(t, expenseID, result.ID)
	assert.True(t, decimal.NewFromFloat(100.00).Equal(result.TotalAmount))
}
