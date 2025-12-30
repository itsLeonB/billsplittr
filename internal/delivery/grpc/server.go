package grpc

import (
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/delivery/grpc/server"
	"github.com/itsLeonB/billsplittr/internal/pkg/config"
	"github.com/itsLeonB/billsplittr/internal/pkg/logger"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/itsLeonB/gerpc"
	"google.golang.org/grpc"
)

func Setup() (*gerpc.GrpcServer, error) {
	providers, err := provider.All()
	if err != nil {
		return nil, err
	}

	servers := server.ProvideServers(providers.Services)

	// Middlewares/Interceptors
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			gerpc.NewLoggingInterceptor(logger.Global),
			gerpc.NewErrorInterceptor(logger.Global),
		),
		grpc.MaxRecvMsgSize(appconstant.MaxFileSize),
	}

	srv := gerpc.NewGrpcServer().
		WithLogger(logger.Global).
		WithAddress(":" + config.Global.App.Port).
		WithOpts(opts...).
		WithRegisterSrvFunc(servers.Register).
		WithShutdownFunc(providers.Shutdown).
		WithReflections()

	return srv, nil
}
