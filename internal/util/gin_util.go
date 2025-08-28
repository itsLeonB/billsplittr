package util

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/ginkgo"
)

func GetProfileID(ctx *gin.Context) (uuid.UUID, error) {
	return ginkgo.GetAndParseFromContext[uuid.UUID](ctx, appconstant.ContextProfileID)
}
