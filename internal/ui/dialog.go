package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func renderDialog(m Model) string {
	width := m.width - 8
	if width < 40 {
		width = 40
	}
	if width > 56 {
		width = 56
	}

	rule := DividerStyle.Render(repeat(CharRule, width))

	title := TitleStyle.Render("Quit TunePipe?")
	prompt := TextDimStyle.Render("Are you sure you want to exit?")

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		KeyHintKeyStyle.Render("Y"),
		TextStyle.Render(" Yes    "),
		KeyHintKeyStyle.Render("N"),
		TextStyle.Render(" No"),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		rule,
		"",
		title,
		prompt,
		"",
		buttons,
		"",
		rule,
	)

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(content))
}
