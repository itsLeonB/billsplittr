package appconstant

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	ErrDataSelect = "error retrieving data"
	ErrDataInsert = "error inserting new data"
	ErrDataUpdate = "error updating data"
	ErrDataDelete = "error deleting data"

	ErrAuthUserNotFound       = "user is not found"
	ErrAuthDuplicateUser      = "user with email %s is already registered"
	ErrAuthUnknownCredentials = "unknown credentials, please check your email/password"

	ErrUserNotFound = "user with ID: %s is not found"
	ErrUserDeleted  = "user with ID: %s is deleted"

	ErrFriendshipNotFound = "friendship not found"
	ErrFriendshipDeleted  = "friendship is deleted"

	ErrTransferMethodNotFound = "transfer method with ID: %s is not found"
	ErrTransferMethodDeleted  = "transfer method with ID: %s is deleted"

	ErrAmountMismatched = "amount mismatch, please check the total amount and the items/fees provided"
	ErrAmountZero       = "amount must be greater than zero"

	ErrNotFriends = "you are not friends with this user, please add them as a friend first"
)

func ErrGroupExpenseNotFound(id uuid.UUID) string {
	return fmt.Sprintf("group expense with ID: %s is not found", id.String())
}

func ErrGroupExpenseDeleted(id uuid.UUID) string {
	return fmt.Sprintf("group expense with ID: %s is deleted", id.String())
}
