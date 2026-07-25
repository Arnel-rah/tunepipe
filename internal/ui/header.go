package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderHeader(m Model, width int) string {
	logo := LogoStyle.Render(IconMusic+" TunePipe")

	stats := fmt.Sprintf("%s %dh%s%s %d",
		StatLabelStyle.Render("listened"),
		m.totalMinutes/60,
		DividerDot(),
		StatLabelStyle.Render("tracks"),
		m.scrobbles,
	)

	gap := width - lipgloss.Width(logo) - lipgloss.Width(stats)
	if gap < 1 {
		gap = 1
	}

	content := fmt.Sprintf("%s%s%s", logo, strings.Repeat(" ", gap), stats)

	return HeaderStyle.Width(width).Render(content)
}
