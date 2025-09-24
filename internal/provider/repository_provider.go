package provider

import (
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
)

type Repositories struct {
	Transactor         crud.Transactor
	GroupExpense       repository.GroupExpenseRepository
	ExpenseItem        repository.ExpenseItemRepository
	ExpenseParticipant repository.ExpenseParticipantRepository
	OtherFee           repository.OtherFeeRepository
	ExpenseBill        repository.ExpenseBillRepository
	TaskQueue          repository.TaskQueue
}

func ProvideRepositories(dbs *DBs, logger ezutil.Logger) *Repositories {
	if dbs.GormDB == nil {
		panic("gormDB cannot be nil")
	}

	return &Repositories{
		Transactor:         crud.NewTransactor(dbs.GormDB),
		GroupExpense:       repository.NewGroupExpenseRepository(dbs.GormDB),
		ExpenseItem:        repository.NewExpenseItemRepository(dbs.GormDB),
		ExpenseParticipant: crud.NewRepository[entity.ExpenseParticipant](dbs.GormDB),
		OtherFee:           repository.NewOtherFeeRepository(dbs.GormDB),
		ExpenseBill:        crud.NewRepository[entity.ExpenseBill](dbs.GormDB),
		TaskQueue:          repository.NewTaskQueue(logger, dbs.Asynq),
	}
}

func (r *Repositories) Shutdown() error {
	return nil
}
