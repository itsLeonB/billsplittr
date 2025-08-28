package provider

import (
	"fmt"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/cocoon-protos/gen/go/auth/v1"
	"github.com/itsLeonB/cocoon-protos/gen/go/friendship/v1"
	"github.com/itsLeonB/cocoon-protos/gen/go/profile/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	conn       *grpc.ClientConn
	Auth       auth.AuthServiceClient
	Profile    profile.ProfileServiceClient
	Friendship friendship.FriendshipServiceClient
}

func ProvideClients(configs config.ServiceClient) *Clients {
	conn, err := grpc.NewClient(
		configs.CocoonHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(fmt.Sprintf("error connecting to grpc client: %v", err))
	}

	return &Clients{
		conn,
		auth.NewAuthServiceClient(conn),
		profile.NewProfileServiceClient(conn),
		friendship.NewFriendshipServiceClient(conn),
	}
}

func (c *Clients) Shutdown() error {
	return c.conn.Close()
}
