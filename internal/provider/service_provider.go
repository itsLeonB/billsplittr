package provider

import (
	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/rotisserie/eris"
)

type Services struct {
	GroupExpense service.GroupExpenseService
	ExpenseItem  service.ExpenseItemService
	OtherFee     service.OtherFeeService
	ExpenseBill  service.ExpenseBillService
}

func ProvideServices(
	repositories *Repositories,
	logger ezutil.Logger,
	cfg config.Config,
	queues *Queues,
) (*Services, error) {
	if repositories == nil {
		return nil, eris.New("repositories cannot be nil")
	}
	if queues == nil {
		return nil, eris.New("queues cannot be nil")
	}

	groupExpenseService := service.NewGroupExpenseService(
		repositories.Transactor,
		repositories.GroupExpense,
		repositories.OtherFee,
		repositories.ExpenseBill,
		service.NewLLMService(cfg.LLM),
		queues.ExpenseBillTextExtracted,
		logger,
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

	expenseBillService := service.NewExpenseBillService(
		repositories.Transactor,
		repositories.ExpenseBill,
		logger,
		queues.OrphanedBillCleanup,
		cfg.BucketNameExpenseBill,
	)

	return &Services{
		GroupExpense: groupExpenseService,
		ExpenseItem:  expenseItemService,
		OtherFee:     otherFeeService,
		ExpenseBill:  expenseBillService,
	}, nil
}
