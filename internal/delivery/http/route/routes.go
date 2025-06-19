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

	apiRoutes := router.Group("/api", errorMiddleware)

	v1 := apiRoutes.Group("/v1")

	authRoutes := v1.Group("/auth")
	authRoutes.POST("/register", handlers.Auth.HandleRegister())
	authRoutes.POST("/login", handlers.Auth.HandleLogin())

	protectedRoutes := v1.Group("/", authMiddleware)

	protectedRoutes.GET("/profile", handlers.Profile.HandleProfile())

	protectedRoutes.POST("/friendships", handlers.Friendship.HandleCreateAnonymousFriendship())
	protectedRoutes.GET("/friendships", handlers.Friendship.HandleGetAll())
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
			return false, nil, ezutil.UnauthorizedError(appconstant.ErrAuthUserNotFound)
		}

		authData := map[string]any{
			appconstant.ContextUserID: userID,
		}

		return true, authData, nil
	}
}
