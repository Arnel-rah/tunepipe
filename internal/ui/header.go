package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderHeader(m Model, width int) string {
	logo := LogoStyle.Render("◆ TUNEPIPE")

	var stats string
	if width < 50 {
		stats = fmt.Sprintf("%dh · %d tracks", m.totalMinutes/60, m.scrobbles)
	} else {
		stats = fmt.Sprintf("%s %dh  %s %d",
			StatLabelStyle.Render("TIME"),
			m.totalMinutes/60,
			StatLabelStyle.Render("PLAYED"),
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
