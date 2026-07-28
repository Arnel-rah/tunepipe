package ui

import "strings"

func renderQueue(m Model, width int) string {
	if len(m.queue) == 0 {
		return TextDimStyle.Render("queue is empty")
	}

	maxItems := len(m.queue)
	if maxItems > 8 {
		maxItems = 8
	}

	start := 0
	if len(m.queue) > maxItems && m.queueCursor >= maxItems {
		start = m.queueCursor - maxItems + 1
		if start > len(m.queue)-maxItems {
			start = len(m.queue) - maxItems
		}
	}
	end := start + maxItems
	if end > len(m.queue) {
		end = len(m.queue)
	}

	lines := make([]string, 0, maxItems+1)
	lines = append(lines, renderTableHeader(width))
	for i := start; i < end; i++ {
		track := m.queue[i]
		lines = append(lines, renderResultLine(track, i == m.queueCursor, width))
	}

	return strings.Join(lines, "\n")
}
