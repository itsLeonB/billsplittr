package provider

import (
	"log"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/itsLeonB/ezutil"
)

type Repositories struct {
	Transactor         ezutil.Transactor
	DebtTransaction    repository.DebtTransactionRepository
	TransferMethod     repository.TransferMethodRepository
	GroupExpense       repository.GroupExpenseRepository
	ExpenseItem        repository.ExpenseItemRepository
	ExpenseParticipant repository.ExpenseParticipantRepository
	OtherFee           repository.OtherFeeRepository
	ExpenseBill        repository.ExpenseBillRepository
	Image              repository.ImageRepository
}

func ProvideRepositories(configs *ezutil.Config) *Repositories {
	googleConfig, ok := configs.Generic.(*config.Google)
	if !ok {
		log.Fatalf("error asserting generic config to google config")
	}

	return &Repositories{
		Transactor:         ezutil.NewTransactor(configs.GORM),
		DebtTransaction:    repository.NewDebtTransactionRepository(configs.GORM),
		TransferMethod:     ezutil.NewCRUDRepository[entity.TransferMethod](configs.GORM),
		GroupExpense:       repository.NewGroupExpenseRepository(configs.GORM),
		ExpenseItem:        repository.NewExpenseItemRepository(configs.GORM),
		ExpenseParticipant: ezutil.NewCRUDRepository[entity.ExpenseParticipant](configs.GORM),
		OtherFee:           repository.NewOtherFeeRepository(configs.GORM),
		ExpenseBill:        ezutil.NewCRUDRepository[entity.ExpenseBill](configs.GORM),
		Image:              repository.NewImageRepository("billsplittr-bills", googleConfig.ServiceAccount),
	}
}
