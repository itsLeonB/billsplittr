package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(
	authService service.AuthService,
) *AuthHandler {
	return &AuthHandler{
		authService,
	}
}

func (ah *AuthHandler) HandleRegister() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request, err := ezutil.BindRequest[dto.RegisterRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		err = ah.authService.Register(ctx, request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusCreated,
			ezutil.NewResponse(appconstant.MsgRegisterSuccess),
		)
	}
}

func (ah *AuthHandler) HandleLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request, err := ezutil.BindRequest[dto.LoginRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		response, err := ah.authService.Login(ctx, request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusOK,
			ezutil.NewResponse(appconstant.MsgLoginSuccess).WithData(response),
		)
	}
}
