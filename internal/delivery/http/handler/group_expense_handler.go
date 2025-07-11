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
	"github.com/itsLeonB/ezutil"
)

type GroupExpenseHandler struct {
	groupExpenseService service.GroupExpenseService
}

func NewGroupExpenseHandler(
	groupExpenseService service.GroupExpenseService,
) *GroupExpenseHandler {
	return &GroupExpenseHandler{
		groupExpenseService,
	}
}

func (geh *GroupExpenseHandler) HandleCreateDraft() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request, err := ezutil.BindRequest[dto.NewGroupExpenseRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		userID, err := util.GetUserID(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request.CreatedByUserID = userID

		response, err := geh.groupExpenseService.CreateDraft(ctx, request)
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

func (geh *GroupExpenseHandler) HandleGetAllCreated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := util.GetUserID(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		groupExpenses, err := geh.groupExpenseService.GetAllCreated(ctx, userID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusOK,
			ezutil.NewResponse(appconstant.MsgGetData).WithData(groupExpenses),
		)
	}
}

func (geh *GroupExpenseHandler) HandleGetDetails() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupExpenseID, err := ezutil.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextGroupExpenseID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		response, err := geh.groupExpenseService.GetDetails(ctx, groupExpenseID)
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
