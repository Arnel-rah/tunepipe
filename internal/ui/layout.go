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

	width := m.width - 4
	if width < 40 {
		width = 40
	}

	if m.confirmQuit {
		return renderDialog(m)
	}

	if len(m.searchResults) == 0 && m.currentTrack == nil && !m.isSearching && !m.searchPending {
		return renderWelcome(m, width)
	}

	header := renderHeader(m, width)
	player := renderPlayer(m, width)
	search := renderSearch(m, width)

	var results string
	if m.showQueue {
		results = renderQueue(m, width)
	} else {
		results = renderResults(m, width)
	}

	progress := renderProgress(m, width)
	footer := renderFooter(m, width)

	var sections []string
	sections = append(sections, header)
	sections = append(sections, "")
	sections = append(sections, player)
	sections = append(sections, "")
	sections = append(sections, search)

	if results != "" {
		sections = append(sections, "")
		sections = append(sections, results)
	}

	if progress != "" {
		sections = append(sections, "")
		sections = append(sections, progress)
	}

	sections = append(sections, "")
	sections = append(sections, footer)

	body := lipgloss.JoinVertical(lipgloss.Left, sections...)

	return lipgloss.NewStyle().
		Padding(0, 1).
		Render(body)
}
