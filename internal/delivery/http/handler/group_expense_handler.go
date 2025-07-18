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

func (geh *GroupExpenseHandler) HandleGetItemDetails() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupExpenseID, err := ezutil.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextGroupExpenseID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		expenseItemID, err := ezutil.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextExpenseItemID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		response, err := geh.groupExpenseService.GetItemDetails(ctx, groupExpenseID, expenseItemID)
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

func (geh *GroupExpenseHandler) HandleUpdateItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupExpenseID, err := ezutil.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextGroupExpenseID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		expenseItemID, err := ezutil.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextExpenseItemID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request, err := ezutil.BindRequest[dto.UpdateExpenseItemRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request.GroupExpenseID = groupExpenseID
		request.ID = expenseItemID

		response, err := geh.groupExpenseService.UpdateItem(ctx, request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusOK,
			ezutil.NewResponse(appconstant.MsgUpdateData).WithData(response),
		)
	}
}

func (geh *GroupExpenseHandler) HandleConfirmDraft() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupExpenseID, err := ezutil.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextGroupExpenseID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		response, err := geh.groupExpenseService.ConfirmDraft(ctx, groupExpenseID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusOK,
			ezutil.NewResponse(appconstant.MsgUpdateData).WithData(response),
		)
	}
}

func (geh *GroupExpenseHandler) HandleGetFeeCalculationMethods() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response := geh.groupExpenseService.GetFeeCalculationMethods()

		ctx.JSON(
			http.StatusOK,
			ezutil.NewResponse(appconstant.MsgGetData).WithData(response),
		)
	}
}

func (geh *GroupExpenseHandler) HandleUpdateFee() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupExpenseID, err := ezutil.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextGroupExpenseID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		otherFeeID, err := ezutil.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextOtherFeeID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request, err := ezutil.BindRequest[dto.UpdateOtherFeeRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request.GroupExpenseID = groupExpenseID
		request.ID = otherFeeID

		response, err := geh.groupExpenseService.UpdateFee(ctx, request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusOK,
			ezutil.NewResponse(appconstant.MsgUpdateData).WithData(response),
		)
	}
}

func (geh *GroupExpenseHandler) HandleAddItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupExpenseID, err := ezutil.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextGroupExpenseID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request, err := ezutil.BindRequest[dto.NewExpenseItemRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request.GroupExpenseID = groupExpenseID

		response, err := geh.groupExpenseService.AddItem(ctx, request)
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

func (geh *GroupExpenseHandler) HandleAddFee() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupExpenseID, err := ezutil.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextGroupExpenseID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request, err := ezutil.BindRequest[dto.NewOtherFeeRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request.GroupExpenseID = groupExpenseID

		response, err := geh.groupExpenseService.AddFee(ctx, request)
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
