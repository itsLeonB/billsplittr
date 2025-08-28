package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/billsplittr/internal/util"
	"github.com/itsLeonB/ginkgo"
)

type DebtHandler struct {
	debtService service.DebtService
}

func NewDebtHandler(debtService service.DebtService) *DebtHandler {
	return &DebtHandler{debtService}
}

func (dh *DebtHandler) HandleCreate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		profileID, err := util.GetProfileID(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request, err := ginkgo.BindRequest[dto.NewDebtTransactionRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request.UserProfileID = profileID

		response, err := dh.debtService.RecordNewTransaction(ctx, request)
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

func (dh *DebtHandler) HandleGetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		profileID, err := util.GetProfileID(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		response, err := dh.debtService.GetTransactions(ctx, profileID)
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
