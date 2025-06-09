package route

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
)

func SetupRoutes(router *gin.Engine, handlers *provider.Handlers, services *provider.Services) {
	// Middlewares
	tokenCheckFunc := newTokenCheckFunc(services.JWT, services.User)
	authMiddleware := ezutil.NewAuthMiddleware("Bearer", tokenCheckFunc)
	errorMiddleware := ezutil.NewErrorMiddleware()

	routeConfigs := []ezutil.RouteConfig{
		{
			Group: "/api",
			Handlers: []gin.HandlerFunc{
				errorMiddleware,
			},
			Versions: []ezutil.RouteVersionConfig{
				{
					Version:  1,
					Handlers: []gin.HandlerFunc{},
					Groups: []ezutil.RouteGroupConfig{
						{
							Group:    "/auth",
							Handlers: []gin.HandlerFunc{},
							Endpoints: []ezutil.EndpointConfig{
								{
									Method:   "POST",
									Endpoint: "/register",
									Handlers: []gin.HandlerFunc{
										handlers.Auth.HandleRegister(),
									},
								},
								{
									Method:   "POST",
									Endpoint: "/login",
									Handlers: []gin.HandlerFunc{
										handlers.Auth.HandleLogin(),
									},
								},
								{
									Method:   "GET",
									Endpoint: "/me",
									Handlers: []gin.HandlerFunc{
										authMiddleware,
										handlers.Auth.HandleProfile(),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	ezutil.SetupRoutes(router, routeConfigs)
}

func newTokenCheckFunc(jwtService ezutil.JWTService, userService service.UserService) func(ctx *gin.Context, token string) (bool, map[string]any, error) {
	return func(ctx *gin.Context, token string) (bool, map[string]any, error) {
		claims, err := jwtService.VerifyToken(token)
		if err != nil {
			return false, nil, err
		}

		tokenUserId, exists := claims.Data[appconstant.ContextUserID]
		if !exists {
			return false, nil, eris.New("missing user ID from token")
		}
		stringUserID, ok := tokenUserId.(string)
		if !ok {
			return false, nil, eris.New("error asserting userID, is not a string")
		}
		userID, err := ezutil.Parse[uuid.UUID](stringUserID)
		if err != nil {
			return false, nil, err
		}

		exists, err = userService.ExistsByID(ctx, userID)
		if err != nil {
			return false, nil, err
		}
		if !exists {
			return false, nil, ezutil.UnauthorizedError(appconstant.MsgAuthUserNotFound)
		}

		authData := map[string]any{
			appconstant.ContextUserID: userID,
		}

		return true, authData, nil
	}
}
