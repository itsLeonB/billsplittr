package provider

import (
	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"gorm.io/gorm"
)

type Repositories struct {
	Transactor         crud.Transactor
	GroupExpense       repository.GroupExpenseRepository
	ExpenseItem        repository.ExpenseItemRepository
	ExpenseParticipant repository.ExpenseParticipantRepository
	OtherFee           repository.OtherFeeRepository
	ExpenseBill        repository.ExpenseBillRepository
}

func ProvideRepositories(gormDB *gorm.DB, googleConfig config.Google, logger ezutil.Logger) *Repositories {
	if gormDB == nil {
		panic("gormDB cannot be nil")
	}

	return &Repositories{
		Transactor:         crud.NewTransactor(gormDB),
		GroupExpense:       repository.NewGroupExpenseRepository(gormDB),
		ExpenseItem:        repository.NewExpenseItemRepository(gormDB),
		ExpenseParticipant: crud.NewCRUDRepository[entity.ExpenseParticipant](gormDB),
		OtherFee:           repository.NewOtherFeeRepository(gormDB),
		ExpenseBill:        crud.NewCRUDRepository[entity.ExpenseBill](gormDB),
	}
}

func (r *Repositories) Shutdown() error {
	return nil
}
