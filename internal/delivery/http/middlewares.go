package http

import (
	"net/http"
	"time"

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

	corsConfig := cors.Config{
		AllowOrigins:     configs.ClientUrls,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Origin", "Cache-Control"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	corsMiddleware := middlewareProvider.NewCorsMiddleware(&corsConfig)
	loggerMiddleware := newLoggerMiddleware(logger)

	return &middlewares{authMiddleware, errorMiddleware, corsMiddleware, loggerMiddleware}
}

// loggerMiddlewareProvider allows customization of the logger interceptor
type loggerMiddlewareProvider struct {
	logger ezutil.Logger
}

// newLoggerMiddleware creates a Gin middleware that logs HTTP requests
// similar to your gRPC interceptor pattern
func newLoggerMiddleware(logger ezutil.Logger) gin.HandlerFunc {
	if logger == nil {
		panic("logger cannot be nil")
	}

	p := loggerMiddlewareProvider{logger}
	return p.handle()
}

func (l *loggerMiddlewareProvider) handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodOptions {
			ctx.Next()
			return
		}

		start := time.Now()
		path := ctx.Request.URL.Path
		method := ctx.Request.Method

		// Build full path with query string for logging
		fullPath := path
		if rawQuery := ctx.Request.URL.RawQuery; rawQuery != "" {
			fullPath = path + "?" + rawQuery
		}

		// Process request
		ctx.Next()

		// Calculate duration
		elapsed := time.Since(start)
		statusCode := ctx.Writer.Status()
		clientIP := ctx.ClientIP()

		// Log based on status code (similar to gRPC error handling)
		if statusCode >= 400 {
			errorMsg := ""
			if len(ctx.Errors) > 0 {
				errorMsg = ctx.Errors.String()
			}

			if errorMsg != "" {
				l.logger.Errorf(
					"[HTTP] method=%s path=%s status=%d duration=%s client_ip=%s error=%s",
					method,
					fullPath,
					statusCode,
					elapsed,
					clientIP,
					errorMsg,
				)
			} else {
				l.logger.Errorf(
					"[HTTP] method=%s path=%s status=%d duration=%s client_ip=%s",
					method,
					fullPath,
					statusCode,
					elapsed,
					clientIP,
				)
			}
		} else {
			l.logger.Infof(
				"[HTTP] method=%s path=%s status=%d duration=%s client_ip=%s",
				method,
				fullPath,
				statusCode,
				elapsed,
				clientIP,
			)
		}
	}
}
