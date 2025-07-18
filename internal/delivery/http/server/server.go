package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/billsplittr/internal/delivery/http/route"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/itsLeonB/ezutil"
)

func SetupHTTPServer(configs *ezutil.Config) *http.Server {
	repositories := provider.ProvideRepositories(configs)
	services := provider.ProvideServices(configs, repositories)
	handlers := provider.ProvideHandlers(services)

	gin.SetMode(configs.App.Env)
	r := gin.Default()
	route.SetupRoutes(r, configs, handlers, services)

	return &http.Server{
		Addr:              fmt.Sprintf(":%s", configs.App.Port),
		Handler:           r,
		ReadTimeout:       configs.App.Timeout,
		ReadHeaderTimeout: configs.App.Timeout,
		WriteTimeout:      configs.App.Timeout,
		IdleTimeout:       configs.App.Timeout,
	}
}

func DefaultConfigs() ezutil.Config {
	timeout, _ := time.ParseDuration("10s")
	tokenDuration, _ := time.ParseDuration("24h")
	cookieDuration, _ := time.ParseDuration("24h")
	secretKey, err := ezutil.GenerateRandomString(32)
	if err != nil {
		log.Fatal("error generating secret key: %w", err)
	}

	appConfig := ezutil.App{
		Env:        "debug",
		Port:       "8080",
		Timeout:    timeout,
		ClientUrls: []string{"http://localhost:3000"},
		Timezone:   "Asia/Jakarta",
	}

	authConfig := ezutil.Auth{
		SecretKey:      secretKey,
		TokenDuration:  tokenDuration,
		CookieDuration: cookieDuration,
		Issuer:         "billsplittr",
		URL:            "http://localhost:8000",
	}

	return ezutil.Config{
		App:  &appConfig,
		Auth: &authConfig,
	}
}
