package dto_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestUploadBillRequestFields(t *testing.T) {
	payerID := uuid.New()
	creatorID := uuid.New()
	imageData := []byte("test image data")

	req := dto.UploadBillRequest{
		PayerProfileID:   payerID,
		CreatorProfileID: creatorID,
		ImageData:        imageData,
		ContentType:      "image/jpeg",
		Filename:         "bill.jpg",
		FileSize:         1024,
	}

	assert.Equal(t, payerID, req.PayerProfileID)
	assert.Equal(t, creatorID, req.CreatorProfileID)
	assert.Equal(t, imageData, req.ImageData)
	assert.Equal(t, "image/jpeg", req.ContentType)
	assert.Equal(t, "bill.jpg", req.Filename)
	assert.Equal(t, int64(1024), req.FileSize)
}

func TestUploadBillResponseFields(t *testing.T) {
	billID := uuid.New()

	resp := dto.UploadBillResponse{
		BillID: billID,
	}

	assert.Equal(t, billID, resp.BillID)
}

func TestExpenseBillResponseFields(t *testing.T) {
	id := uuid.New()
	payerID := uuid.New()
	creatorID := uuid.New()
	groupExpenseID := uuid.New()
	now := time.Now()

	resp := dto.ExpenseBillResponse{
		ID:               id,
		PayerProfileID:   payerID,
		CreatorProfileID: creatorID,
		GroupExpenseID:   groupExpenseID,
		ImageURL:         "https://example.com/bill.jpg",
		Filename:         "bill.jpg",
		ContentType:      "image/jpeg",
		FileSize:         1024,
		Status:           "uploaded",
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	assert.Equal(t, id, resp.ID)
	assert.Equal(t, payerID, resp.PayerProfileID)
	assert.Equal(t, creatorID, resp.CreatorProfileID)
	assert.Equal(t, groupExpenseID, resp.GroupExpenseID)
	assert.Equal(t, "https://example.com/bill.jpg", resp.ImageURL)
	assert.Equal(t, "bill.jpg", resp.Filename)
	assert.Equal(t, "image/jpeg", resp.ContentType)
	assert.Equal(t, int64(1024), resp.FileSize)
	assert.Equal(t, "uploaded", resp.Status)
}
