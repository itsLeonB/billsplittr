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

type UpdateExpenseItemRequest struct {
	ID             uuid.UUID                   `json:"-"`
	GroupExpenseID uuid.UUID                   `json:"-"`
	Name           string                      `json:"name" binding:"required,min=3"`
	Amount         decimal.Decimal             `json:"amount" binding:"required"`
	Quantity       int                         `json:"quantity" binding:"required,min=1"`
	Participants   []NewItemParticipantRequest `json:"participants" binding:"dive"`
}

type NewItemParticipantRequest struct {
	ProfileID uuid.UUID       `json:"profileId" binding:"required"`
	Share     decimal.Decimal `json:"share" binding:"required"`
}

type GroupExpenseResponse struct {
	ID                    uuid.UUID                    `json:"id"`
	PayerProfileID        uuid.UUID                    `json:"payerProfileId"`
	PayerName             string                       `json:"payerName,omitempty"`
	PaidByUser            bool                         `json:"paidByUser"`
	TotalAmount           decimal.Decimal              `json:"totalAmount"`
	Description           string                       `json:"description"`
	Items                 []ExpenseItemResponse        `json:"items,omitempty"`
	OtherFees             []OtherFeeResponse           `json:"otherFees,omitempty"`
	CreatorProfileID      uuid.UUID                    `json:"creatorProfileId"`
	CreatorName           string                       `json:"creatorName,omitempty"`
	CreatedByUser         bool                         `json:"createdByUser"`
	Confirmed             bool                         `json:"confirmed"`
	ParticipantsConfirmed bool                         `json:"participantsConfirmed"`
	CreatedAt             time.Time                    `json:"createdAt"`
	UpdatedAt             time.Time                    `json:"updatedAt"`
	DeletedAt             time.Time                    `json:"deletedAt,omitzero"`
	Participants          []ExpenseParticipantResponse `json:"participants,omitempty"`
}

type ExpenseItemResponse struct {
	ID             uuid.UUID                 `json:"id"`
	GroupExpenseID uuid.UUID                 `json:"groupExpenseId"`
	Name           string                    `json:"name"`
	Amount         decimal.Decimal           `json:"amount"`
	Quantity       int                       `json:"quantity"`
	CreatedAt      time.Time                 `json:"createdAt"`
	UpdatedAt      time.Time                 `json:"updatedAt"`
	DeletedAt      time.Time                 `json:"deletedAt,omitzero"`
	Participants   []ItemParticipantResponse `json:"participants,omitempty"`
}

type OtherFeeResponse struct {
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Amount    decimal.Decimal `json:"amount"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	DeletedAt time.Time       `json:"deletedAt,omitzero"`
}

type ItemParticipantResponse struct {
	ProfileName string          `json:"profileName"`
	ProfileID   uuid.UUID       `json:"profileId"`
	Share       decimal.Decimal `json:"share"`
	IsUser      bool            `json:"isUser"`
}

type ExpenseParticipantResponse struct {
	ProfileName string          `json:"profileName"`
	ProfileID   uuid.UUID       `json:"profileId"`
	ShareAmount decimal.Decimal `json:"shareAmount"`
	IsUser      bool            `json:"isUser"`
}
