package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/itsLeonB/ginkgo"
)

func Setup(configs config.Config) *ginkgo.HttpServer {
	providers := provider.All(configs)

	gin.SetMode(configs.App.Env)
	r := gin.New()
	registerRoutes(r, configs, providers.Logger, providers.Services)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", configs.App.Port),
		Handler:           r,
		ReadTimeout:       configs.App.Timeout,
		ReadHeaderTimeout: configs.App.Timeout,
		WriteTimeout:      configs.App.Timeout,
		IdleTimeout:       configs.App.Timeout,
	}

	shutdownFunc := func() error {
		if err := providers.Clients.Shutdown(); err != nil {
			return err
		}
		return providers.DBs.Shutdown()
	}

	return ginkgo.NewHttpServer(srv, configs.Timeout, providers.Logger, shutdownFunc)
}
