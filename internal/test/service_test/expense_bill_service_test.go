package service_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestExpenseBillService_Creation(t *testing.T) {
	expenseBillService := service.NewExpenseBillService(nil, nil, nil)
	
	assert.NotNil(t, expenseBillService)
}
