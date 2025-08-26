package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/shopspring/decimal"
)

type NewDebtTransactionRequest struct {
	UserProfileID    uuid.UUID
	FriendProfileID  uuid.UUID          `json:"friendProfileId" binding:"required"`
	Action           appconstant.Action `json:"action" binding:"oneof=LEND BORROW RECEIVE RETURN"`
	Amount           decimal.Decimal    `json:"amount" binding:"required"`
	TransferMethodID uuid.UUID          `json:"transferMethodId" binding:"required"`
	Description      string             `json:"description"`
}

type DebtTransactionResponse struct {
	ID             uuid.UUID                       `json:"id"`
	ProfileID      uuid.UUID                       `json:"profileId"`
	Type           appconstant.DebtTransactionType `json:"type"`
	Action         appconstant.Action              `json:"action"`
	Amount         decimal.Decimal                 `json:"amount"`
	TransferMethod string                          `json:"transferMethod"`
	Description    string                          `json:"description"`
	CreatedAt      time.Time                       `json:"createdAt"`
	UpdatedAt      time.Time                       `json:"updatedAt"`
	DeletedAt      time.Time                       `json:"deletedAt,omitzero"`
}
