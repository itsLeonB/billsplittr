package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil"
)

type FriendshipHandler struct {
	friendshipService service.FriendshipService
}

func NewFriendshipHandler(
	friendshipService service.FriendshipService,
) *FriendshipHandler {
	return &FriendshipHandler{
		friendshipService,
	}
}

func (fh *FriendshipHandler) HandleCreateAnonymousFriendship() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := ezutil.GetFromContext[uuid.UUID](ctx, appconstant.ContextUserID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request, err := ezutil.BindRequest[dto.NewAnonymousFriendshipRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request.UserID = userID

		response, err := fh.friendshipService.CreateAnonymous(ctx, request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusCreated,
			ezutil.NewResponse(appconstant.MsgInsertData).WithData(response),
		)
	}
}
