package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/delivery/http/handler"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/itsLeonB/ezutil/v2"
)

func registerRoutes(router *gin.Engine, configs config.Config, logger ezutil.Logger, services *provider.Services) {
	handlers := handler.ProvideHandlers(services)
	middlewares := provideMiddlewares(configs.App, logger, services.Auth)

	router.Use(middlewares.logger, middlewares.cors)

	apiRoutes := router.Group("/api", middlewares.err)

	v1 := apiRoutes.Group("/v1")

	authRoutes := v1.Group("/auth")
	authRoutes.POST("/register", handlers.Auth.HandleRegister())
	authRoutes.POST("/login", handlers.Auth.HandleLogin())

	protectedRoutes := v1.Group("/", middlewares.auth)

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
	groupExpenseRoutes.GET(fmt.Sprintf("/:%s/items/:%s", appconstant.ContextGroupExpenseID, appconstant.ContextExpenseItemID), handlers.GroupExpense.HandleGetItemDetails())
	groupExpenseRoutes.PUT(fmt.Sprintf("/:%s/items/:%s", appconstant.ContextGroupExpenseID, appconstant.ContextExpenseItemID), handlers.GroupExpense.HandleUpdateItem())
	groupExpenseRoutes.PATCH(fmt.Sprintf("/:%s/confirmed", appconstant.ContextGroupExpenseID), handlers.GroupExpense.HandleConfirmDraft())
	groupExpenseRoutes.GET("/fee-calculation-methods", handlers.GroupExpense.HandleGetFeeCalculationMethods())
	groupExpenseRoutes.PUT(fmt.Sprintf("/:%s/fees/:%s", appconstant.ContextGroupExpenseID, appconstant.ContextOtherFeeID), handlers.GroupExpense.HandleUpdateFee())
	groupExpenseRoutes.POST(fmt.Sprintf("/:%s/items", appconstant.ContextGroupExpenseID), handlers.GroupExpense.HandleAddItem())
	groupExpenseRoutes.POST(fmt.Sprintf("/:%s/fees", appconstant.ContextGroupExpenseID), handlers.GroupExpense.HandleAddFee())
	groupExpenseRoutes.DELETE(fmt.Sprintf("/:%s/items/:%s", appconstant.ContextGroupExpenseID, appconstant.ContextExpenseItemID), handlers.GroupExpense.HandleRemoveItem())
	groupExpenseRoutes.DELETE(fmt.Sprintf("/:%s/fees/:%s", appconstant.ContextGroupExpenseID, appconstant.ContextOtherFeeID), handlers.GroupExpense.HandleRemoveFee())
	groupExpenseRoutes.POST("/bills", handlers.GroupExpense.HandleUploadBill())
}
