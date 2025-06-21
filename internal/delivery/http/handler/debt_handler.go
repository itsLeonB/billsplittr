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

type DebtHandler struct {
	debtService service.DebtService
}

func NewDebtHandler(debtService service.DebtService) *DebtHandler {
	return &DebtHandler{debtService}
}

func (dh *DebtHandler) HandleCreate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request, err := ezutil.BindRequest[dto.NewDebtTransactionRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		response, err := dh.debtService.RecordNewTransaction(ctx, request)
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

func (dh *DebtHandler) HandleGetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := ezutil.GetFromContext[uuid.UUID](ctx, appconstant.ContextUserID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		response, err := dh.debtService.GetTransactions(ctx, userID)
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
