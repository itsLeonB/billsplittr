package provider

import (
	"github.com/itsLeonB/billsplittr/internal/service"
)

type Services struct {
	Auth           service.AuthService
	Profile        service.ProfileService
	Friendship     service.FriendshipService
	Debt           service.DebtService
	TransferMethod service.TransferMethodService
	GroupExpense   service.GroupExpenseService
	ExpenseBill    service.ExpenseBillService
}

func ProvideServices(repositories *Repositories, clients *Clients) *Services {
	if repositories == nil {
		panic("repositories cannot be nil")
	}
	if clients == nil {
		panic("clients cannot be nil")
	}

	authService := service.NewAuthService(clients.Auth)

	profileService := service.NewProfileService(clients.Profile)

	friendshipService := service.NewFriendshipService(
		repositories.DebtTransaction,
		clients.Friendship,
	)

	debtService := service.NewDebtService(
		repositories.Transactor,
		repositories.DebtTransaction,
		repositories.TransferMethod,
		repositories.GroupExpense,
		friendshipService,
	)

	transferMethodService := service.NewTransferMethodService(repositories.TransferMethod)

	groupExpenseService := service.NewGroupExpenseService(
		repositories.Transactor,
		repositories.GroupExpense,
		friendshipService,
		repositories.ExpenseItem,
		debtService,
		repositories.OtherFee,
		profileService,
	)

	expenseBillService := service.NewExpenseBillService(
		friendshipService,
		repositories.Image,
		repositories.ExpenseBill,
	)

	return &Services{
		Auth:           authService,
		Profile:        profileService,
		Friendship:     friendshipService,
		Debt:           debtService,
		TransferMethod: transferMethodService,
		GroupExpense:   groupExpenseService,
		ExpenseBill:    expenseBillService,
	}
}
