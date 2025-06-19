package provider

import (
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil"
)

type Services struct {
	Auth       service.AuthService
	User       service.UserService
	JWT        ezutil.JWTService
	Friendship service.FriendshipService
}

func ProvideServices(configs *ezutil.Config, repositories *Repositories) *Services {
	hashService := ezutil.NewHashService(10)
	jwtService := ezutil.NewJwtService(configs.Auth)

	authService := service.NewAuthService(
		hashService,
		jwtService,
		repositories.User,
		repositories.Transactor,
		repositories.UserProfile,
	)

	userService := service.NewUserService(
		repositories.User,
		repositories.UserProfile,
	)

	friendshipService := service.NewFriendshipRepository(
		repositories.Transactor,
		repositories.User,
		repositories.UserProfile,
		repositories.Friendship,
	)

	return &Services{
		Auth:       authService,
		User:       userService,
		JWT:        jwtService,
		Friendship: friendshipService,
	}
}
