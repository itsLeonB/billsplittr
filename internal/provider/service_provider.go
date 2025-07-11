package provider

import (
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil"
)

type Services struct {
	Auth           service.AuthService
	User           service.UserService
	JWT            ezutil.JWTService
	Friendship     service.FriendshipService
	Debt           service.DebtService
	TransferMethod service.TransferMethodService
	GroupExpense   service.GroupExpenseService
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
		repositories.Transactor,
		repositories.User,
		repositories.UserProfile,
	)

	friendshipService := service.NewFriendshipService(
		repositories.Transactor,
		repositories.UserProfile,
		repositories.Friendship,
		userService,
		repositories.DebtTransaction,
	)

	debtService := service.NewDebtService(
		repositories.Transactor,
		repositories.Friendship,
		userService,
		repositories.DebtTransaction,
		repositories.TransferMethod,
	)

	transferMethodService := service.NewTransferMethodService(repositories.TransferMethod)

	groupExpenseService := service.NewGroupExpenseService(
		repositories.Transactor,
		repositories.GroupExpense,
		userService,
		friendshipService,
		repositories.ExpenseItem,
	)

	return &Services{
		Auth:           authService,
		User:           userService,
		JWT:            jwtService,
		Friendship:     friendshipService,
		Debt:           debtService,
		TransferMethod: transferMethodService,
		GroupExpense:   groupExpenseService,
	}
}
