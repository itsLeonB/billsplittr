package provider

import (
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil"
)

type Services struct {
	Auth service.AuthService
	User service.UserService
	JWT  ezutil.JWTService
}

func ProvideServices(configs *ezutil.Config, repositories *Repositories) *Services {
	transactor := ezutil.NewTransactor(configs.GORM)
	hashService := ezutil.NewHashService(10)
	jwtService := ezutil.NewJwtService(configs.Auth)
	authService := service.NewAuthService(hashService, jwtService, repositories.User, transactor)
	userService := service.NewUserService(repositories.User)

	return &Services{
		Auth: authService,
		User: userService,
		JWT:  jwtService,
	}
}
