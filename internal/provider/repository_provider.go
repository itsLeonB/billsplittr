package provider

import (
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/itsLeonB/go-crud"
)

type Repositories struct {
	Transactor   crud.Transactor
	GroupExpense repository.GroupExpenseRepository
	ExpenseItem  repository.ExpenseItemRepository
	OtherFee     repository.OtherFeeRepository
	ExpenseBill  repository.ExpenseBillRepository
}

func ProvideRepositories(dbs *DBs) *Repositories {
	if dbs.GormDB == nil {
		panic("gormDB cannot be nil")
	}

	return &Repositories{
		Transactor:   crud.NewTransactor(dbs.GormDB),
		GroupExpense: repository.NewGroupExpenseRepository(dbs.GormDB),
		ExpenseItem:  repository.NewExpenseItemRepository(dbs.GormDB),
		OtherFee:     repository.NewOtherFeeRepository(dbs.GormDB),
		ExpenseBill:  crud.NewRepository[entity.ExpenseBill](dbs.GormDB),
	}
}
