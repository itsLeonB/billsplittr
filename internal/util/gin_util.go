package util

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
)

func GetUserID(ctx *gin.Context) (uuid.UUID, error) {
	return ezutil.GetFromContext[uuid.UUID](ctx, appconstant.ContextUserID)
}

func FindUserID(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(appconstant.ContextUserID).(uuid.UUID)
	if !ok {
		return uuid.Nil, eris.New("failed to retrieve user ID from context")
	}
	return userID, nil
}
