package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestDebtService_RecordNewTransaction_InvalidAmount(t *testing.T) {
	debtService := service.NewDebtService(nil, nil)

	request := dto.NewDebtTransactionRequest{
		UserProfileID:   uuid.New(),
		FriendProfileID: uuid.New(),
		Amount:          decimal.Zero,
	}

	_, err := debtService.RecordNewTransaction(context.Background(), request)
	
	assert.Error(t, err)
	assert.NotNil(t, err)
}
