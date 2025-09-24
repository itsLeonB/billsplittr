package grpc

import (
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/delivery/grpc/server"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/itsLeonB/gerpc"
	"google.golang.org/grpc"
)

func Setup(configs config.Config) *gerpc.GrpcServer {
	logger := provider.ProvideLogger(config.AppName, configs.Env)
	providers := provider.All(configs, logger)
	servers := server.ProvideServers(providers.Services)

	// Middlewares/Interceptors
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			gerpc.NewLoggingInterceptor(providers.Logger),
			gerpc.NewErrorInterceptor(providers.Logger),
		),
		grpc.MaxRecvMsgSize(appconstant.MaxFileSize),
	}

	return gerpc.NewGrpcServer().
		WithLogger(providers.Logger).
		WithAddress(":" + configs.App.Port).
		WithOpts(opts...).
		WithRegisterSrvFunc(servers.Register).
		WithShutdownFunc(providers.Shutdown)
}
