package route

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
)

func SetupRoutes(router *gin.Engine, configs *ezutil.Config, handlers *provider.Handlers, services *provider.Services) {
	tokenCheckFunc := newTokenCheckFunc(services.JWT, services.User)
	authMiddleware := ezutil.NewAuthMiddleware("Bearer", tokenCheckFunc)
	errorMiddleware := ezutil.NewErrorMiddleware()

	corsConfig := cors.Config{
		AllowOrigins:     configs.App.ClientUrls,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Origin", "Cache-Control"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	corsMiddleware := ezutil.NewCorsMiddleware(&corsConfig)

	// Apply CORS middleware to the entire router first
	router.Use(corsMiddleware)

	apiRoutes := router.Group("/api", errorMiddleware)

	v1 := apiRoutes.Group("/v1")

	authRoutes := v1.Group("/auth")
	authRoutes.POST("/register", handlers.Auth.HandleRegister())
	authRoutes.POST("/login", handlers.Auth.HandleLogin())

	protectedRoutes := v1.Group("/", authMiddleware)

	protectedRoutes.GET("/profile", handlers.Profile.HandleProfile())

	friendshipRoutes := protectedRoutes.Group("/friendships")
	friendshipRoutes.POST("", handlers.Friendship.HandleCreateAnonymousFriendship())
	friendshipRoutes.GET("", handlers.Friendship.HandleGetAll())
	friendshipRoutes.GET(fmt.Sprintf("/:%s", appconstant.ContextFriendshipID), handlers.Friendship.HandleGetDetails())

	protectedRoutes.GET("/transfer-methods", handlers.TransferMethod.HandleGetAll())

	protectedRoutes.POST("/debts", handlers.Debt.HandleCreate())
	protectedRoutes.GET("/debts", handlers.Debt.HandleGetAll())

	groupExpenseRoutes := protectedRoutes.Group("/group-expenses")
	groupExpenseRoutes.POST("", handlers.GroupExpense.HandleCreateDraft())
	groupExpenseRoutes.GET("", handlers.GroupExpense.HandleGetAllCreated())
	groupExpenseRoutes.GET(fmt.Sprintf("/:%s", appconstant.ContextGroupExpenseID), handlers.GroupExpense.HandleGetDetails())
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
