package service_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestTransferMethodService_Creation(t *testing.T) {
	transferMethodService := service.NewTransferMethodService(nil)
	
	assert.NotNil(t, transferMethodService)
}
