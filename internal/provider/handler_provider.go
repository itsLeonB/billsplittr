package provider

import "github.com/itsLeonB/billsplittr/internal/delivery/http/handler"

type Handlers struct {
	Auth *handler.AuthHandler
}

func ProvideHandlers(services *Services) *Handlers {
	return &Handlers{
		Auth: handler.NewAuthHandler(services.Auth, services.User),
	}
}
