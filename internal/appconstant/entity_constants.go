package appconstant

type FeeCalculationMethod string

const (
	EqualSplitFee    FeeCalculationMethod = "EQUAL_SPLIT"
	ItemizedSplitFee FeeCalculationMethod = "ITEMIZED_SPLIT"
)

type ExpenseStatus string

const (
	DraftExpense     ExpenseStatus = "DRAFT"
	ReadyExpense     ExpenseStatus = "READY"
	ConfirmedExpense ExpenseStatus = "CONFIRMED"
)

type BillStatus string

const (
	PendingBill BillStatus = "PENDING"
	ParsedBill  BillStatus = "PARSED"
	FailedBill  BillStatus = "FAILED"
)
