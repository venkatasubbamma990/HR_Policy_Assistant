package query

import (
	"regexp"
	"strings"
)

var spacePattern = regexp.MustCompile(`\s+`)

// Normalize trims and collapses whitespace. Lowercase is used only for intent detection.
func Normalize(question string) string {
	question = strings.TrimSpace(question)
	question = spacePattern.ReplaceAllString(question, " ")
	return question
}

// NormalizedForSearch returns a lowercase copy for keyword-based intent detection.
func NormalizedForSearch(question string) string {
	return strings.ToLower(Normalize(question))
}
