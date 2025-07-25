package provider

import "github.com/itsLeonB/billsplittr/internal/delivery/http/handler"

type Handlers struct {
	Auth           *handler.AuthHandler
	Friendship     *handler.FriendshipHandler
	Profile        *handler.ProfileHandler
	TransferMethod *handler.TransferMethodHandler
	Debt           *handler.DebtHandler
	GroupExpense   *handler.GroupExpenseHandler
}

func ProvideHandlers(services *Services) *Handlers {
	return &Handlers{
		Auth:           handler.NewAuthHandler(services.Auth),
		Friendship:     handler.NewFriendshipHandler(services.Friendship),
		Profile:        handler.NewProfileHandler(services.User),
		TransferMethod: handler.NewTransferMethodHandler(services.TransferMethod),
		Debt:           handler.NewDebtHandler(services.Debt),
		GroupExpense:   handler.NewGroupExpenseHandler(services.GroupExpense, services.ExpenseBill),
	}
}
