package util

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/ezutil"
)

func GetUserID(ctx *gin.Context) (uuid.UUID, error) {
	return ezutil.GetFromContext[uuid.UUID](ctx, appconstant.ContextUserID)
}
