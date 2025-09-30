package appconstant

import "errors"

const (
	ErrDataUpdate = "error updating data"

	ErrAmountMismatched = "amount mismatch, please check the total amount and the items/fees provided"
	ErrAmountZero       = "amount must be greater than zero"

	ErrProcessFile = "error processing file upload"

	ErrNonPositiveAmount = "amount must be positive (>0)"

	ErrInvalidFileType     = "invalid file type"
	ErrFileTooLarge        = "file too large"
	ErrProfileNotFound     = "profile not found"
	ErrStorageUploadFailed = "storage upload failed"
	ErrBillNotFound        = "bill not found"
	ErrUnauthorized        = "unauthorized access"

	ErrStructValidation = "failed to validate struct"

	ErrNilRequest = "request is nil"
)

var ErrExpenseNotDetected = errors.New("NOT_DETECTED")
