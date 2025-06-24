package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil"
)

type ProfileHandler struct {
	userService service.UserService
}

func NewProfileHandler(
	userService service.UserService,
) *ProfileHandler {
	return &ProfileHandler{
		userService,
	}
}

func (ph *ProfileHandler) HandleProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		parsedUserID, err := ezutil.GetFromContext[uuid.UUID](ctx, appconstant.ContextUserID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		response, err := ph.userService.GetProfile(ctx, parsedUserID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusOK,
			ezutil.NewResponse(appconstant.MsgGetData).WithData(response),
		)
	}
}
