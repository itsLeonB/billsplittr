package util

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/ezutil"
)

func GetNameFromEmail(email string) string {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) < 2 || parts[0] == "" {
		return ""
	}
	localPart := parts[0]

	re := regexp.MustCompile(`[a-zA-Z]+`)
	matches := re.FindAllString(localPart, -1)
	if len(matches) > 0 {
		name := matches[0]
		return ezutil.Capitalize(name)
	}

	return ""
}

func NotFoundMessage(ent entity.Entity) string {
	return fmt.Sprintf("%s with ID: %s is not found", ent.SimpleName(), ent.GetID())
}

func DeletedMessage(ent entity.Entity) string {
	return fmt.Sprintf("%s with ID: %s is deleted", ent.SimpleName(), ent.GetID())
}
