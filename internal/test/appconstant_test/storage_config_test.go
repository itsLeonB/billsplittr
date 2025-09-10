package appconstant_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/stretchr/testify/assert"
)

func TestStorageConfigConstants(t *testing.T) {
	assert.Equal(t, int64(10*1024*1024), int64(appconstant.MaxFileSize))
	assert.Equal(t, "pending", appconstant.StatusPending)
	assert.Equal(t, "uploaded", appconstant.StatusUploaded)
	assert.Equal(t, "processed", appconstant.StatusProcessed)
	assert.Equal(t, "failed", appconstant.StatusFailed)
}

func TestAllowedContentTypesMap(t *testing.T) {
	expectedTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/webp"}
	
	assert.Len(t, appconstant.AllowedContentTypes, 4)
	
	for _, contentType := range expectedTypes {
		_, exists := appconstant.AllowedContentTypes[contentType]
		assert.True(t, exists, "Content type %s should be allowed", contentType)
	}
	
	// Test non-allowed type
	_, exists := appconstant.AllowedContentTypes["image/gif"]
	assert.False(t, exists, "Content type image/gif should not be allowed")
}
