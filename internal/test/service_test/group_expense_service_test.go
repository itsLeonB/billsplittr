package service_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestGroupExpenseService_GetFeeCalculationMethods(t *testing.T) {
	groupExpenseService := service.NewGroupExpenseService(
		nil, nil, nil, nil, nil, nil, nil,
	)

	methods := groupExpenseService.GetFeeCalculationMethods()

	assert.NotNil(t, methods)
	assert.GreaterOrEqual(t, len(methods), 0)
}
