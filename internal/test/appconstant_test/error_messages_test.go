package appconstant_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/stretchr/testify/assert"
)

func TestErrorMessagesConstants(t *testing.T) {
	assert.Equal(t, "error updating data", appconstant.ErrDataUpdate)
	assert.Equal(t, "amount mismatch, please check the total amount and the items/fees provided", appconstant.ErrAmountMismatched)
	assert.Equal(t, "amount must be greater than zero", appconstant.ErrAmountZero)
	assert.Equal(t, "error processing file upload", appconstant.ErrProcessFile)
	assert.Equal(t, "amount must be positive (>0)", appconstant.ErrNonPositiveAmount)
	assert.Equal(t, "invalid file type", appconstant.ErrInvalidFileType)
	assert.Equal(t, "file too large", appconstant.ErrFileTooLarge)
	assert.Equal(t, "profile not found", appconstant.ErrProfileNotFound)
	assert.Equal(t, "storage upload failed", appconstant.ErrStorageUploadFailed)
	assert.Equal(t, "bill not found", appconstant.ErrBillNotFound)
	assert.Equal(t, "unauthorized access", appconstant.ErrUnauthorized)
	assert.Equal(t, "failed to validate struct", appconstant.ErrStructValidation)
}
