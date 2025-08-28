package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ginkgo"
)

type TransferMethodHandler struct {
	transferMethodService service.TransferMethodService
}

func NewTransferMethodHandler(transferMethodService service.TransferMethodService) *TransferMethodHandler {
	return &TransferMethodHandler{transferMethodService}
}

func (tmh *TransferMethodHandler) HandleGetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response, err := tmh.transferMethodService.GetAll(ctx)
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
