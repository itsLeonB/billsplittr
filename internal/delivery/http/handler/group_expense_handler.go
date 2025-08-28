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
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/ginkgo"
	"github.com/itsLeonB/ungerr"
	"github.com/rotisserie/eris"
)

type GroupExpenseHandler struct {
	groupExpenseService service.GroupExpenseService
	expenseBillService  service.ExpenseBillService
}

func NewGroupExpenseHandler(
	groupExpenseService service.GroupExpenseService,
	expenseBillService service.ExpenseBillService,
) *GroupExpenseHandler {
	return &GroupExpenseHandler{
		groupExpenseService,
		expenseBillService,
	}
}

func (geh *GroupExpenseHandler) HandleCreateDraft() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request, err := ginkgo.BindRequest[dto.NewGroupExpenseRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		profileID, err := util.GetProfileID(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request.CreatedByProfileID = profileID

		response, err := geh.groupExpenseService.CreateDraft(ctx, request)
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

func (geh *GroupExpenseHandler) HandleGetAllCreated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		profileID, err := util.GetProfileID(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		groupExpenses, err := geh.groupExpenseService.GetAllCreated(ctx, profileID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusOK,
			ginkgo.NewResponse(appconstant.MsgGetData).WithData(groupExpenses),
		)
	}
}

func (geh *GroupExpenseHandler) HandleGetDetails() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupExpenseID, err := ginkgo.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextGroupExpenseID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		profileID, err := util.GetProfileID(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		response, err := geh.groupExpenseService.GetDetails(ctx, groupExpenseID, profileID)
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

func (geh *GroupExpenseHandler) HandleGetItemDetails() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupExpenseID, err := ginkgo.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextGroupExpenseID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		expenseItemID, err := ginkgo.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextExpenseItemID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		profileID, err := util.GetProfileID(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		response, err := geh.groupExpenseService.GetItemDetails(ctx, groupExpenseID, expenseItemID, profileID)
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

func (geh *GroupExpenseHandler) HandleUpdateItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupExpenseID, err := ginkgo.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextGroupExpenseID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		expenseItemID, err := ginkgo.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextExpenseItemID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		profileID, err := util.GetProfileID(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request, err := ginkgo.BindRequest[dto.UpdateExpenseItemRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request.GroupExpenseID = groupExpenseID
		request.ID = expenseItemID

		response, err := geh.groupExpenseService.UpdateItem(ctx, profileID, request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusOK,
			ginkgo.NewResponse(appconstant.MsgUpdateData).WithData(response),
		)
	}
}

func (geh *GroupExpenseHandler) HandleConfirmDraft() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupExpenseID, err := ginkgo.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextGroupExpenseID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		profileID, err := util.GetProfileID(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		response, err := geh.groupExpenseService.ConfirmDraft(ctx, groupExpenseID, profileID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusOK,
			ginkgo.NewResponse(appconstant.MsgUpdateData).WithData(response),
		)
	}
}

func (geh *GroupExpenseHandler) HandleGetFeeCalculationMethods() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response := geh.groupExpenseService.GetFeeCalculationMethods()

		ctx.JSON(
			http.StatusOK,
			ginkgo.NewResponse(appconstant.MsgGetData).WithData(response),
		)
	}
}

func (geh *GroupExpenseHandler) HandleUpdateFee() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupExpenseID, err := ginkgo.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextGroupExpenseID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		otherFeeID, err := ginkgo.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextOtherFeeID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		profileID, err := util.GetProfileID(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request, err := ginkgo.BindRequest[dto.UpdateOtherFeeRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request.GroupExpenseID = groupExpenseID
		request.ID = otherFeeID

		response, err := geh.groupExpenseService.UpdateFee(ctx, profileID, request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusOK,
			ginkgo.NewResponse(appconstant.MsgUpdateData).WithData(response),
		)
	}
}

func (geh *GroupExpenseHandler) HandleAddItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupExpenseID, err := ginkgo.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextGroupExpenseID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		profileID, err := util.GetProfileID(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request, err := ginkgo.BindRequest[dto.NewExpenseItemRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request.GroupExpenseID = groupExpenseID

		response, err := geh.groupExpenseService.AddItem(ctx, profileID, request)
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

func (geh *GroupExpenseHandler) HandleAddFee() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupExpenseID, err := ginkgo.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextGroupExpenseID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		profileID, err := util.GetProfileID(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request, err := ginkgo.BindRequest[dto.NewOtherFeeRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request.GroupExpenseID = groupExpenseID

		response, err := geh.groupExpenseService.AddFee(ctx, profileID, request)
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

func (geh *GroupExpenseHandler) HandleRemoveItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupExpenseID, err := ginkgo.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextGroupExpenseID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		expenseItemID, err := ginkgo.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextExpenseItemID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request := dto.DeleteExpenseItemRequest{
			ID:             expenseItemID,
			GroupExpenseID: groupExpenseID,
		}

		if err = geh.groupExpenseService.RemoveItem(ctx, request); err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

func (geh *GroupExpenseHandler) HandleRemoveFee() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupExpenseID, err := ginkgo.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextGroupExpenseID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		feeID, err := ginkgo.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextOtherFeeID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request := dto.DeleteOtherFeeRequest{
			ID:             feeID,
			GroupExpenseID: groupExpenseID,
		}

		if err = geh.groupExpenseService.RemoveFee(ctx, request); err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

func (geh *GroupExpenseHandler) HandleUploadBill() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		profileID, err := util.GetProfileID(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		payerProfileID := uuid.Nil
		if payerProfileIDStr := ctx.PostForm("payerProfileId"); payerProfileIDStr != "" {
			payerProfileID, err = ezutil.Parse[uuid.UUID](payerProfileIDStr)
			if err != nil {
				_ = ctx.Error(err)
				return
			}
		}

		fileHeader, err := ctx.FormFile("bill")
		if err != nil {
			_ = ctx.Error(eris.Wrap(err, appconstant.ErrProcessFile))
			return
		}

		contentType, ok := util.IsImageType(fileHeader)
		if !ok {
			_ = ctx.Error(ungerr.BadRequestError("file is not an image"))
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			_ = ctx.Error(eris.Wrap(err, appconstant.ErrProcessFile))
			return
		}

		request := dto.NewExpenseBillRequest{
			PayerProfileID:   payerProfileID,
			CreatorProfileID: profileID,
			ImageReader:      file,
			ContentType:      contentType,
		}

		if err = geh.expenseBillService.Upload(ctx, request); err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, ginkgo.NewResponse(appconstant.MsgBillUploaded))
	}
}
