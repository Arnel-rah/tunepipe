package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderHeader(m Model, width int) string {
	brand := lipgloss.JoinHorizontal(
		lipgloss.Left,
		LogoStyle.Render("◆ TUNEPIPE"),
		"  ",
		HeaderSubtitleStyle.Render("AUDIO PLAYER"),
	)

	var stats string
	if width < 50 {
		stats = fmt.Sprintf("%dh %d", m.totalMinutes/60, m.scrobbles)
	} else if width < 70 {
		stats = fmt.Sprintf("%s %dh  %s %d",
			HeaderSubtitleStyle.Render("TIME"),
			m.totalMinutes/60,
			HeaderSubtitleStyle.Render("PLAYED"),
			m.scrobbles,
		)
	} else {
		stats = fmt.Sprintf("%s %s  %s  %s %s",
			StatLabelStyle.Render("TIME"),
			HeaderStatValueStyle.Render(fmt.Sprintf("%dh", m.totalMinutes/60)),
			HeaderStatSeparatorStyle.Render("·"),
			StatLabelStyle.Render("PLAYED"),
			HeaderStatValueStyle.Render(fmt.Sprintf("%d", m.scrobbles)),
		)
	}

	if width < 50 {
		brand = LogoStyle.Render("◆ TUNEPIPE")
	}

	gap := width - lipgloss.Width(brand) - lipgloss.Width(stats)
	if gap < 1 {
		gap = 1
	}

	content := fmt.Sprintf("%s%s%s", brand, strings.Repeat(" ", gap), stats)

	return HeaderStyle.Width(width).Render(content)
}
