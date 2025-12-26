package uuidutil

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil/v2"
)

func Parse(id string) (uuid.UUID, error) {
	return ezutil.Parse[uuid.UUID](id)
}
