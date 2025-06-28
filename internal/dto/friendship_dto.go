package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/shopspring/decimal"
)

type NewAnonymousFriendshipRequest struct {
	UserID uuid.UUID
	Name   string `json:"name" binding:"required,min=3"`
}

type FriendshipResponse struct {
	ID          uuid.UUID                  `json:"id"`
	Type        appconstant.FriendshipType `json:"type"`
	ProfileID   uuid.UUID                  `json:"profileId"`
	ProfileName string                     `json:"profileName"`
	CreatedAt   time.Time                  `json:"createdAt"`
	UpdatedAt   time.Time                  `json:"updatedAt"`
	DeletedAt   time.Time                  `json:"deletedAt,omitzero"`
}

type FriendshipWithProfile struct {
	Friendship    FriendshipResponse
	UserProfile   ProfileResponse
	FriendProfile ProfileResponse
}

type FriendDetails struct {
	ID        uuid.UUID                  `json:"id"`
	ProfileID uuid.UUID                  `json:"profileId"`
	Name      string                     `json:"name"`
	Type      appconstant.FriendshipType `json:"type"`
	Email     string                     `json:"email,omitempty"`
	Phone     string                     `json:"phone,omitempty"`
	Avatar    string                     `json:"avatar,omitempty"`
	CreatedAt time.Time                  `json:"createdAt"`
	UpdatedAt time.Time                  `json:"updatedAt"`
	DeletedAt time.Time                  `json:"deletedAt,omitzero"`
}

type FriendBalance struct {
	TotalOwedToYou decimal.Decimal      `json:"totalOwedToYou"`
	TotalYouOwe    decimal.Decimal      `json:"totalYouOwe"`
	NetBalance     decimal.Decimal      `json:"netBalance"`
	Currency       appconstant.Currency `json:"currency"`
}

type FriendStats struct {
	TotalTransactions        int             `json:"totalTransactions"`
	FirstTransactionDate     time.Time       `json:"firstTransactionDate"`
	LastTransactionDate      time.Time       `json:"lastTransactionDate"`
	MostUsedTransferMethod   string          `json:"mostUsedTransferMethod"`
	AverageTransactionAmount decimal.Decimal `json:"averageTransactionAmount"`
}

type FriendDetailsResponse struct {
	Friend       FriendDetails             `json:"friend"`
	Balance      FriendBalance             `json:"balance"`
	Transactions []DebtTransactionResponse `json:"transactions"`
	Stats        FriendStats               `json:"stats"`
}
