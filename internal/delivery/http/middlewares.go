package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/ginkgo"
)

type middlewares struct {
	auth   gin.HandlerFunc
	err    gin.HandlerFunc
	cors   gin.HandlerFunc
	logger gin.HandlerFunc
}

func provideMiddlewares(configs config.App, logger ezutil.Logger, authSvc service.AuthService) *middlewares {
	tokenCheckFunc := func(ctx *gin.Context, token string) (bool, map[string]any, error) {
		return authSvc.VerifyToken(ctx, token)
	}

	middlewareProvider := ginkgo.NewMiddlewareProvider(logger)
	authMiddleware := middlewareProvider.NewAuthMiddleware("Bearer", tokenCheckFunc)
	errorMiddleware := middlewareProvider.NewErrorMiddleware()
	loggingMiddleware := middlewareProvider.NewLoggingMiddleware()

	corsMiddleware := middlewareProvider.NewCorsMiddleware(&cors.Config{
		AllowOrigins:     configs.ClientUrls,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Origin", "Cache-Control"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})

	return &middlewares{authMiddleware, errorMiddleware, corsMiddleware, loggingMiddleware}
}
