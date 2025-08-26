package provider

import (
	"os"

	"github.com/itsLeonB/billsplittr/internal/logging"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/cocoon-protos/gen/go/auth"
	"github.com/itsLeonB/cocoon-protos/gen/go/friendship"
	"github.com/itsLeonB/cocoon-protos/gen/go/profile"
	"github.com/itsLeonB/ezutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func ProvideServices(configs *ezutil.Config, repositories *Repositories) *Services {
	cocoonHost := os.Getenv("COCOON_HOST")

	conn, err := grpc.NewClient(
		cocoonHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logging.Logger.Fatalf("error connecting to grpc client: %v", err)
	}

	authClient := auth.NewAuthServiceClient(conn)
	profileClient := profile.NewProfileServiceClient(conn)
	friendshipClient := friendship.NewFriendshipServiceClient(conn)

	authService := service.NewAuthService(authClient)

	profileService := service.NewProfileService(profileClient)

	friendshipService := service.NewFriendshipService(
		repositories.DebtTransaction,
		friendshipClient,
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
