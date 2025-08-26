package util

import (
	"fmt"

	"github.com/itsLeonB/billsplittr/internal/entity"
)

func NotFoundMessage(ent entity.Entity) string {
	return fmt.Sprintf("%s with ID: %s is not found", ent.SimpleName(), ent.GetID())
}

func DeletedMessage(ent entity.Entity) string {
	return fmt.Sprintf("%s with ID: %s is deleted", ent.SimpleName(), ent.GetID())
}
