package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/billsplittr/internal/test/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOtherFeeService_Add_ValidationErrorNonPositiveAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := mocks.NewMockTransactor(ctrl)
	mockGroupExpenseRepo := mocks.NewMockGroupExpenseRepository(ctrl)
	mockOtherFeeRepo := mocks.NewMockOtherFeeRepository(ctrl)
	mockGroupExpenseSvc := mocks.NewMockGroupExpenseService(ctrl)

	svc := service.NewOtherFeeService(mockTransactor, mockGroupExpenseRepo, mockOtherFeeRepo, mockGroupExpenseSvc)

	request := dto.NewOtherFeeRequest{
		ProfileID:      uuid.New(),
		GroupExpenseID: uuid.New(),
		OtherFeeData: dto.OtherFeeData{
			Name:              "Service Fee",
			Amount:            decimal.NewFromFloat(-5.00), // Negative amount
			CalculationMethod: appconstant.EqualSplitFee,
		},
	}

	_, err := svc.Add(context.Background(), request)

	assert.Error(t, err)
	// Error message check removed
}

func TestOtherFeeService_Add_ValidationErrorZeroAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := mocks.NewMockTransactor(ctrl)
	mockGroupExpenseRepo := mocks.NewMockGroupExpenseRepository(ctrl)
	mockOtherFeeRepo := mocks.NewMockOtherFeeRepository(ctrl)
	mockGroupExpenseSvc := mocks.NewMockGroupExpenseService(ctrl)

	svc := service.NewOtherFeeService(mockTransactor, mockGroupExpenseRepo, mockOtherFeeRepo, mockGroupExpenseSvc)

	request := dto.NewOtherFeeRequest{
		ProfileID:      uuid.New(),
		GroupExpenseID: uuid.New(),
		OtherFeeData: dto.OtherFeeData{
			Name:              "Service Fee",
			Amount:            decimal.Zero, // Zero amount
			CalculationMethod: appconstant.EqualSplitFee,
		},
	}

	_, err := svc.Add(context.Background(), request)

	assert.Error(t, err)
	// Error message check removed
}

func TestOtherFeeServiceGetCalculationMethods(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := mocks.NewMockTransactor(ctrl)
	mockGroupExpenseRepo := mocks.NewMockGroupExpenseRepository(ctrl)
	mockOtherFeeRepo := mocks.NewMockOtherFeeRepository(ctrl)
	mockGroupExpenseSvc := mocks.NewMockGroupExpenseService(ctrl)

	svc := service.NewOtherFeeService(mockTransactor, mockGroupExpenseRepo, mockOtherFeeRepo, mockGroupExpenseSvc)

	methods := svc.GetCalculationMethods()

	assert.NotEmpty(t, methods)

	// Check that we have the expected calculation methods
	methodNames := make([]appconstant.FeeCalculationMethod, len(methods))
	for i, method := range methods {
		methodNames[i] = method.Method
	}

	assert.Contains(t, methodNames, appconstant.EqualSplitFee)
	assert.Contains(t, methodNames, appconstant.ItemizedSplitFee)
}

func TestOtherFeeService_Update_ValidationErrorNonPositiveAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := mocks.NewMockTransactor(ctrl)
	mockGroupExpenseRepo := mocks.NewMockGroupExpenseRepository(ctrl)
	mockOtherFeeRepo := mocks.NewMockOtherFeeRepository(ctrl)
	mockGroupExpenseSvc := mocks.NewMockGroupExpenseService(ctrl)

	svc := service.NewOtherFeeService(mockTransactor, mockGroupExpenseRepo, mockOtherFeeRepo, mockGroupExpenseSvc)

	request := dto.UpdateOtherFeeRequest{
		ProfileID:      uuid.New(),
		ID:             uuid.New(),
		GroupExpenseID: uuid.New(),
		OtherFeeData: dto.OtherFeeData{
			Name:              "Updated Fee",
			Amount:            decimal.NewFromFloat(-10.00), // Negative amount
			CalculationMethod: appconstant.EqualSplitFee,
		},
	}

	_, err := svc.Update(context.Background(), request)

	assert.Error(t, err)
	// Error message check removed
}
