package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/billsplittr/internal/util"
	"github.com/itsLeonB/ginkgo"
)

type FriendshipHandler struct {
	friendshipService service.FriendshipService
	friendDetailsSvc  service.FriendDetailsService
}

func NewFriendshipHandler(
	friendshipService service.FriendshipService,
	friendDetailsSvc service.FriendDetailsService,
) *FriendshipHandler {
	return &FriendshipHandler{
		friendshipService,
		friendDetailsSvc,
	}
}

func (fh *FriendshipHandler) HandleCreateAnonymousFriendship() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		profileID, err := util.GetProfileID(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request, err := ginkgo.BindRequest[dto.NewAnonymousFriendshipRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request.ProfileID = profileID

		response, err := fh.friendshipService.CreateAnonymous(ctx, request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusCreated,
			ginkgo.NewResponse(appconstant.MsgInsertData).WithData(response),
		)
	}
}

func (fh *FriendshipHandler) HandleGetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		profileID, err := util.GetProfileID(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		response, err := fh.friendshipService.GetAll(ctx, profileID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusOK,
			ginkgo.NewResponse(appconstant.MsgGetData).WithData(response),
		)
	}
}

func (fh *FriendshipHandler) HandleGetDetails() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		profileID, err := util.GetProfileID(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		friendshipID, err := ginkgo.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextFriendshipID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		response, err := fh.friendDetailsSvc.GetDetails(ctx, profileID, friendshipID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusOK,
			ginkgo.NewResponse(appconstant.MsgGetData).WithData(response),
		)
	}
}
