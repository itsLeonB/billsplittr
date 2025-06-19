package appconstant

type DebtTransactionType string
type FriendshipType string

const (
	Lend  DebtTransactionType = "LEND"
	Repay DebtTransactionType = "REPAY"

	Real      FriendshipType = "REAL"
	Anonymous FriendshipType = "ANON"
)
