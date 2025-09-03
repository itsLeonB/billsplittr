package provider

import (
	"github.com/itsLeonB/billsplittr/internal/service"
)

type Services struct {
	Auth           service.AuthService
	Profile        service.ProfileService
	Friendship     service.FriendshipService
	TransferMethod service.TransferMethodService
	Debt           service.DebtService
	FriendDetails  service.FriendDetailsService
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
		clients.Friendship,
	)

	transferMethodService := service.NewTransferMethodService(clients.TransferMethod)

	debtService := service.NewDebtService(
		clients.Debt,
		friendshipService,
	)

	friendDetailsService := service.NewFriendDetailsService(clients.Friendship, debtService)

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
		TransferMethod: transferMethodService,
		Debt:           debtService,
		FriendDetails:  friendDetailsService,
		GroupExpense:   groupExpenseService,
		ExpenseBill:    expenseBillService,
	}
}
