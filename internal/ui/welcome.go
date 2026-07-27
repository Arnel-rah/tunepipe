package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func renderWelcome(m Model, width int) string {
	if len(m.searchResults) > 0 || m.currentTrack != nil {
		return ""
	}

	height := m.height - 4
	if height < 10 {
		height = 10
	}

	title := lipgloss.NewStyle().
		Foreground(PrimaryPink).
		Bold(true).
		Render("♪ TunePipe")

	subtitle := lipgloss.NewStyle().
		Foreground(TextDimColor).
		Italic(true).
		Render("Terminal audio player")

	searchHint := lipgloss.NewStyle().
		Foreground(Amber).
		Bold(true).
		Render("Press / to start searching")

	type shortcut struct {
		key  string
		desc string
	}
	shortcuts := []shortcut{
		{"/", "Search music"},
		{"j/k", "Navigate results"},
		{"Enter", "Play selected"},
		{"Space", "Pause/Resume"},
		{"q", "Quit"},
	}

	keyColWidth := 0
	for _, s := range shortcuts {
		if l := lipgloss.Width(s.key); l > keyColWidth {
			keyColWidth = l
		}
	}

	var shortcutsLines []string
	for _, s := range shortcuts {
		key := KeyHintKeyStyle.
			Width(keyColWidth).
			Align(lipgloss.Right).
			Render(s.key)
		desc := KeyHintTextStyle.Render(s.desc)
		shortcutsLines = append(shortcutsLines, fmt.Sprintf("%s  %s", key, desc))
	}
	shortcutsBlock := lipgloss.JoinVertical(lipgloss.Left, shortcutsLines...)

	shortcutsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(TextDimColor).
		Padding(1, 3).
		Render(shortcutsBlock)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(title),
		lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(subtitle),
		"",
		lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(searchHint),
		"",
		lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(shortcutsBox),
	)

	return lipgloss.NewStyle().
		Height(height).
		Width(width).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)
}
