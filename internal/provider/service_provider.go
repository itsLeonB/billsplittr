package provider

import (
	"github.com/itsLeonB/billsplittr/internal/client"
	"github.com/itsLeonB/billsplittr/internal/pkg/config"
	"github.com/itsLeonB/billsplittr/internal/service"
)

type Services struct {
	GroupExpense service.GroupExpenseService
	ExpenseItem  service.ExpenseItemService
	OtherFee     service.OtherFeeService
	ExpenseBill  service.ExpenseBillService
}

func ProvideServices(
	repositories *Repositories,
	queues *Queues,
) (*Services, error) {
	groupExpenseService := service.NewGroupExpenseService(
		repositories.Transactor,
		repositories.GroupExpense,
		repositories.OtherFee,
		repositories.ExpenseBill,
		service.NewLLMService(),
	)

	expenseItemService := service.NewExpenseItemService(
		repositories.Transactor,
		repositories.GroupExpense,
		repositories.ExpenseItem,
		groupExpenseService,
	)

	otherFeeService := service.NewOtherFeeService(
		repositories.Transactor,
		repositories.GroupExpense,
		repositories.OtherFee,
		groupExpenseService,
	)

	ocr, err := client.NewOCRClient()
	if err != nil {
		return nil, err
	}

	expenseBillService := service.NewExpenseBillService(
		repositories.Transactor,
		repositories.ExpenseBill,
		queues.OrphanedBillCleanup,
		config.Global.BucketNameExpenseBill,
		ocr,
	)

	return &Services{
		GroupExpense: groupExpenseService,
		ExpenseItem:  expenseItemService,
		OtherFee:     otherFeeService,
		ExpenseBill:  expenseBillService,
	}, nil
}
