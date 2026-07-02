package chunk

import "strings"

// EstimateTokens approximates token count (~4 characters per token for English prose).
func EstimateTokens(text string) int {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}
	return (len(text) + 3) / 4
}

// OverlapText returns the trailing portion of text that fits within targetTokens.
func OverlapText(text string, targetTokens int) string {
	text = strings.TrimSpace(text)
	if text == "" || targetTokens <= 0 {
		return ""
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var selected []string
	tokens := 0
	for i := len(words) - 1; i >= 0; i-- {
		wordTokens := EstimateTokens(words[i])
		if tokens+wordTokens > targetTokens && len(selected) > 0 {
			break
		}
		selected = append([]string{words[i]}, selected...)
		tokens += wordTokens
		if tokens >= targetTokens {
			break
		}
	}

	return strings.Join(selected, " ")
}
