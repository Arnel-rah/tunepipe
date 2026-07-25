package ui

import (
	"unicode/utf8"
)

func truncate(s string, maxLen int) string {
	if maxLen < 1 {
		return ""
	}

	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}

	runes := []rune(s)
	if maxLen < 3 {
		return string(runes[:maxLen])
	}

	return string(runes[:maxLen-3]) + "…"
}
