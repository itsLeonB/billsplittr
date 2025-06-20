package appconstant

const (
	ErrDataSelect = "error retrieving data"
	ErrDataInsert = "error inserting new data"
	ErrDataUpdate = "error updating data"
	ErrDataDelete = "error deleting data"

	ErrAuthUserNotFound       = "user is not found"
	ErrAuthDuplicateUser      = "user with username %s is already registered"
	ErrAuthUnknownCredentials = "unknown credentials, please check your username and password"

	ErrUserNotFound = "user with ID: %s is not found"
	ErrUserDeleted  = "user with ID: %s is deleted"

	ErrFriendshipNotFound = "friendship not found"
	ErrFriendshipDeleted  = "friendship is deleted"

	ErrTransferMethodNotFound = "transfer method with ID: %s is not found"
	ErrTransferMethodDeleted  = "transfer method with ID: %s is deleted"
)
