package util

import (
	"regexp"
	"strings"
	"unicode"
)

func GetNameFromEmail(email string) string {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) == 0 {
		return ""
	}
	localPart := parts[0]

	re := regexp.MustCompile(`[a-zA-Z]+`)
	matches := re.FindAllString(localPart, -1)
	if len(matches) > 0 {
		name := matches[0]
		return capitalize(name)
	}

	return ""
}

func capitalize(word string) string {
	if len(word) == 0 {
		return ""
	}
	return string(unicode.ToUpper(rune(word[0]))) + strings.ToLower(word[1:])
}
