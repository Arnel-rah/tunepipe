package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderHeader(m Model, width int) string {
	logo := LogoStyle.Render(IconMusic + " TunePipe")

	var stats string
	if width < 50 {
		stats = fmt.Sprintf("%dh %d", m.totalMinutes/60, m.scrobbles)
	} else {
		stats = fmt.Sprintf("%s %dh %s %d",
			StatLabelStyle.Render("listened"),
			m.totalMinutes/60,
			StatLabelStyle.Render("tracks"),
			m.scrobbles,
		)
	}

	gap := width - lipgloss.Width(logo) - lipgloss.Width(stats)
	if gap < 1 {
		gap = 1
	}

	content := fmt.Sprintf("%s%s%s", logo, strings.Repeat(" ", gap), stats)

	return HeaderStyle.Width(width).Render(content)
}
