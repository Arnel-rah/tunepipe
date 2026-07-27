package ui

import (
	"strings"
)

func renderQueue(m Model, width int) string {
	if len(m.queue) == 0 {
		return TextDimStyle.Render("queue is empty")
	}

	lines := make([]string, 0, len(m.queue))
	for i, track := range m.queue {
		lines = append(lines, renderResultLine(track, i == m.queueCursor, true, false, width))
	}

	return strings.Join(lines, "\n")
}
