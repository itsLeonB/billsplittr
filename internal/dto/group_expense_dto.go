package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type NewGroupExpenseRequest struct {
	PayerProfileID     uuid.UUID               `json:"payerProfileId"`
	TotalAmount        decimal.Decimal         `json:"totalAmount" binding:"required"`
	Description        string                  `json:"description"`
	Items              []NewExpenseItemRequest `json:"items" binding:"required,min=1,dive"`
	OtherFees          []NewOtherFeeRequest    `json:"otherFees" binding:"dive"`
	CreatedByUserID    uuid.UUID               `json:"-"`
	CreatedByProfileID uuid.UUID               `json:"-"`
}

type NewExpenseItemRequest struct {
	Name     string          `json:"name" binding:"required,min=3"`
	Amount   decimal.Decimal `json:"amount" binding:"required"`
	Quantity int             `json:"quantity" binding:"required,min=1"`
}

type NewOtherFeeRequest struct {
	Name   string          `json:"name" binding:"required,min=3"`
	Amount decimal.Decimal `json:"amount" binding:"required"`
}

type GroupExpenseResponse struct {
	ID                 uuid.UUID             `json:"id"`
	PayerProfileID     uuid.UUID             `json:"payerProfileId"`
	TotalAmount        decimal.Decimal       `json:"totalAmount"`
	Description        string                `json:"description"`
	Items              []ExpenseItemResponse `json:"items,omitempty"`
	OtherFees          []OtherFeeResponse    `json:"otherFees,omitempty"`
	CreatedByProfileID uuid.UUID             `json:"createdByProfileId"`
	CreatedAt          time.Time             `json:"createdAt"`
	UpdatedAt          time.Time             `json:"updatedAt"`
	DeletedAt          time.Time             `json:"deletedAt,omitzero"`
}

type ExpenseItemResponse struct {
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Amount    decimal.Decimal `json:"amount"`
	Quantity  int             `json:"quantity"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	DeletedAt time.Time       `json:"deletedAt,omitzero"`
}

type OtherFeeResponse struct {
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Amount    decimal.Decimal `json:"amount"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	DeletedAt time.Time       `json:"deletedAt,omitzero"`
}
