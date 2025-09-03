package provider

import (
	"errors"
	"fmt"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/cocoon-protos/gen/go/auth/v1"
	"github.com/itsLeonB/cocoon-protos/gen/go/friendship/v1"
	"github.com/itsLeonB/cocoon-protos/gen/go/profile/v1"
	"github.com/itsLeonB/drex-protos/gen/go/debt/v1"
	"github.com/itsLeonB/drex-protos/gen/go/transaction/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	conns          []*grpc.ClientConn
	Auth           auth.AuthServiceClient
	Profile        profile.ProfileServiceClient
	Friendship     friendship.FriendshipServiceClient
	TransferMethod transaction.TransferMethodServiceClient
	Debt           debt.DebtServiceClient
}

func ProvideClients(configs config.ServiceClient) *Clients {
	cocoonConn, err := grpc.NewClient(
		configs.CocoonHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(fmt.Sprintf("error connecting to grpc client: %v", err))
	}

	drexConn, err := grpc.NewClient(
		configs.DrexHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(fmt.Sprintf("error connecting to grpc client: %v", err))
	}

	return &Clients{
		[]*grpc.ClientConn{cocoonConn, drexConn},
		auth.NewAuthServiceClient(cocoonConn),
		profile.NewProfileServiceClient(cocoonConn),
		friendship.NewFriendshipServiceClient(cocoonConn),
		transaction.NewTransferMethodServiceClient(drexConn),
		debt.NewDebtServiceClient(drexConn),
	}
}

func (c *Clients) Shutdown() error {
	var err error
	for _, conn := range c.conns {
		if e := conn.Close(); e != nil {
			err = errors.Join(err, e)
		}
	}
	return err
}
