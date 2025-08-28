package handler

import (
	"github.com/itsLeonB/billsplittr/internal/provider"
)

type Handlers struct {
	Auth           *AuthHandler
	Friendship     *FriendshipHandler
	Profile        *ProfileHandler
	TransferMethod *TransferMethodHandler
	Debt           *DebtHandler
	GroupExpense   *GroupExpenseHandler
}

func ProvideHandlers(services *provider.Services) *Handlers {
	return &Handlers{
		Auth:           NewAuthHandler(services.Auth),
		Friendship:     NewFriendshipHandler(services.Friendship),
		Profile:        NewProfileHandler(services.Profile),
		TransferMethod: NewTransferMethodHandler(services.TransferMethod),
		Debt:           NewDebtHandler(services.Debt),
		GroupExpense:   NewGroupExpenseHandler(services.GroupExpense, services.ExpenseBill),
	}
}
