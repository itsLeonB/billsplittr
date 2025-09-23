package provider

import (
	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil/v2"
)

type Services struct {
	GroupExpense service.GroupExpenseService
	ExpenseItem  service.ExpenseItemService
	OtherFee     service.OtherFeeService
	ExpenseBill  service.ExpenseBillService
}

func ProvideServices(googleConfig config.Google, repositories *Repositories, logger ezutil.Logger) *Services {
	if repositories == nil {
		panic("repositories cannot be nil")
	}

	groupExpenseService := service.NewGroupExpenseService(
		repositories.Transactor,
		repositories.GroupExpense,
		repositories.OtherFee,
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
	)

	return &Services{
		GroupExpense: groupExpenseService,
		ExpenseItem:  expenseItemService,
		OtherFee:     otherFeeService,
		ExpenseBill:  expenseBillService,
	}
}
