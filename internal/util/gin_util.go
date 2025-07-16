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

func GetProfileID(ctx context.Context) (uuid.UUID, error) {
	profileID, ok := ctx.Value(appconstant.ContextProfileID).(uuid.UUID)
	if !ok {
		return uuid.Nil, eris.New("failed to retrieve profile ID from context")
	}
	return profileID, nil
}
