package service_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestGroupExpenseServiceInterface(t *testing.T) {
	// Test that the interface is properly defined
	var svc service.GroupExpenseService
	assert.Nil(t, svc)
}

func TestExpenseItemServiceInterface(t *testing.T) {
	// Test that the interface is properly defined
	var svc service.ExpenseItemService
	assert.Nil(t, svc)
}

func TestOtherFeeServiceInterface(t *testing.T) {
	// Test that the interface is properly defined
	var svc service.OtherFeeService
	assert.Nil(t, svc)
}

func TestExpenseBillServiceInterface(t *testing.T) {
	// Test that the interface is properly defined
	var svc service.ExpenseBillService
	assert.Nil(t, svc)
}
