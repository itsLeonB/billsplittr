package entity_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestStorageUploadRequestFields(t *testing.T) {
	req := entity.StorageUploadRequest{
		Data:        []byte("test data"),
		ContentType: "image/jpeg",
		Filename:    "test.jpg",
		BucketName:  "test-bucket",
		ObjectKey:   "test-key",
	}

	assert.Equal(t, []byte("test data"), req.Data)
	assert.Equal(t, "image/jpeg", req.ContentType)
	assert.Equal(t, "test.jpg", req.Filename)
	assert.Equal(t, "test-bucket", req.BucketName)
	assert.Equal(t, "test-key", req.ObjectKey)
}

func TestStorageUploadResponseFields(t *testing.T) {
	resp := entity.StorageUploadResponse{
		URL:       "https://example.com/test.jpg",
		ObjectKey: "test-key",
	}

	assert.Equal(t, "https://example.com/test.jpg", resp.URL)
	assert.Equal(t, "test-key", resp.ObjectKey)
}
