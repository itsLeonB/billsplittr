package dto

import (
	"time"

	"github.com/google/uuid"
)

type UploadBillRequest struct {
	PayerProfileID   uuid.UUID `validate:"required"`
	CreatorProfileID uuid.UUID `validate:"required"`
	ImageData        []byte    `validate:"required"`
	ContentType      string    `validate:"required,oneof=image/jpeg image/png image/jpg image/webp"`
	Filename         string    `validate:"required"`
	FileSize         int64     `validate:"required"`
}

type UploadBillResponse struct {
	BillID uuid.UUID
}

type ExpenseBillResponse struct {
	ID               uuid.UUID
	PayerProfileID   uuid.UUID
	CreatorProfileID uuid.UUID
	GroupExpenseID   uuid.UUID
	ImageURL         string
	Filename         string
	ContentType      string
	FileSize         int64
	Status           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time
}
