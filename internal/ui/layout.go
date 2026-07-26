package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.width == 0 {
		m.width = 80
	}
	if m.height == 0 {
		m.height = 24
	}

	width := m.width - 6
	if width < 40 {
		width = 40
	}

	if m.confirmQuit {
		return renderDialog(m)
	}

	header := renderHeader(m, width)
	player := renderPlayer(m, width)
	search := renderSearch(m, width)
	results := renderResults(m, width)
	progress := renderProgress(m, width)
	footer := renderFooter(m, width)

	sections := []string{header, player, search}
	if results != "" {
		sections = append(sections, results)
	}
	if progress != "" {
		sections = append(sections, progress)
	}

	rule := DividerStyle.Render(repeat(CharRule, width))

	rows := []string{sections[0]}
	for _, s := range sections[1:] {
		rows = append(rows, "", rule, "", s)
	}
	rows = append(rows, "", rule, "", footer)

	body := lipgloss.JoinVertical(lipgloss.Left, rows...)

	return lipgloss.NewStyle().
		Padding(1, 3).
		Render(body)
}
