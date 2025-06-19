package provider

import "github.com/itsLeonB/billsplittr/internal/delivery/http/handler"

type Handlers struct {
	Auth       *handler.AuthHandler
	Friendship *handler.FriendshipHandler
	Profile    *handler.ProfileHandler
}

func ProvideHandlers(services *Services) *Handlers {
	return &Handlers{
		Auth:       handler.NewAuthHandler(services.Auth),
		Friendship: handler.NewFriendshipHandler(services.Friendship),
		Profile:    handler.NewProfileHandler(services.User),
	}
}
