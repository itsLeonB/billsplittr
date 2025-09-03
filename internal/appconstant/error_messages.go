package appconstant

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

	ErrAmountMismatched = "amount mismatch, please check the total amount and the items/fees provided"
	ErrAmountZero       = "amount must be greater than zero"

	ErrNotFriends = "you are not friends with this user, please add them as a friend first"

	ErrProcessFile = "error processing file upload"

	ErrNonPositiveAmount = "amount must be positive (>0)"

	ErrServiceClient = "service client communication failure"
)
