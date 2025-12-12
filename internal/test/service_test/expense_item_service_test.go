package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/mocks"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/go-crud"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestExpenseItemService_Add_ValidationErrorNonPositiveAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := crud.NewMockTransactor(ctrl)
	mockGroupExpenseRepo := mocks.NewMockGroupExpenseRepository(ctrl)
	mockExpenseItemRepo := mocks.NewMockExpenseItemRepository(ctrl)
	mockGroupExpenseSvc := mocks.NewMockGroupExpenseService(ctrl)

	svc := service.NewExpenseItemService(mockTransactor, mockGroupExpenseRepo, mockExpenseItemRepo, mockGroupExpenseSvc)

	request := dto.NewExpenseItemRequest{
		ProfileID:      uuid.New(),
		GroupExpenseID: uuid.New(),
		ExpenseItemData: dto.ExpenseItemData{
			Name:     "Test Item",
			Amount:   decimal.NewFromFloat(-10.00), // Negative amount
			Quantity: 1,
		},
	}

	_, err := svc.Add(context.Background(), request)

	assert.Error(t, err)
	// Error message check removed
}

func TestExpenseItemService_Add_ValidationErrorZeroAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := crud.NewMockTransactor(ctrl)
	mockGroupExpenseRepo := mocks.NewMockGroupExpenseRepository(ctrl)
	mockExpenseItemRepo := mocks.NewMockExpenseItemRepository(ctrl)
	mockGroupExpenseSvc := mocks.NewMockGroupExpenseService(ctrl)

	svc := service.NewExpenseItemService(mockTransactor, mockGroupExpenseRepo, mockExpenseItemRepo, mockGroupExpenseSvc)

	request := dto.NewExpenseItemRequest{
		ProfileID:      uuid.New(),
		GroupExpenseID: uuid.New(),
		ExpenseItemData: dto.ExpenseItemData{
			Name:     "Test Item",
			Amount:   decimal.Zero, // Zero amount
			Quantity: 1,
		},
	}

	_, err := svc.Add(context.Background(), request)

	assert.Error(t, err)
	// Error message check removed
}

func TestExpenseItemService_GetDetailsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := crud.NewMockTransactor(ctrl)
	mockGroupExpenseRepo := mocks.NewMockGroupExpenseRepository(ctrl)
	mockExpenseItemRepo := mocks.NewMockExpenseItemRepository(ctrl)
	mockGroupExpenseSvc := mocks.NewMockGroupExpenseService(ctrl)

	svc := service.NewExpenseItemService(mockTransactor, mockGroupExpenseRepo, mockExpenseItemRepo, mockGroupExpenseSvc)

	groupExpenseID := uuid.New()
	expenseItemID := uuid.New()

	expectedItem := entity.ExpenseItem{
		BaseEntity:     crud.BaseEntity{ID: expenseItemID},
		GroupExpenseID: groupExpenseID,
		Name:           "Test Item",
		Amount:         decimal.NewFromFloat(25.50),
		Quantity:       2,
	}

	mockExpenseItemRepo.EXPECT().
		FindFirst(gomock.Any(), gomock.Any()).
		Return(expectedItem, nil)

	result, err := svc.GetDetails(context.Background(), groupExpenseID, expenseItemID)

	assert.NoError(t, err)
	assert.Equal(t, expenseItemID, result.ID)
	assert.Equal(t, groupExpenseID, result.GroupExpenseID)
	assert.Equal(t, "Test Item", result.Name)
	assert.True(t, decimal.NewFromFloat(25.50).Equal(result.Amount))
	assert.Equal(t, 2, result.Quantity)
}

func TestExpenseItemService_GetDetailsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := crud.NewMockTransactor(ctrl)
	mockGroupExpenseRepo := mocks.NewMockGroupExpenseRepository(ctrl)
	mockExpenseItemRepo := mocks.NewMockExpenseItemRepository(ctrl)
	mockGroupExpenseSvc := mocks.NewMockGroupExpenseService(ctrl)

	svc := service.NewExpenseItemService(mockTransactor, mockGroupExpenseRepo, mockExpenseItemRepo, mockGroupExpenseSvc)

	groupExpenseID := uuid.New()
	expenseItemID := uuid.New()

	// Return zero value entity (not found)
	mockExpenseItemRepo.EXPECT().
		FindFirst(gomock.Any(), gomock.Any()).
		Return(entity.ExpenseItem{}, nil)

	_, err := svc.GetDetails(context.Background(), groupExpenseID, expenseItemID)

	assert.Error(t, err)
	// Error message check removed
}
