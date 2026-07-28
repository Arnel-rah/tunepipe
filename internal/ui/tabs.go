package ui

import "strings"

func renderTabs(m Model, width int) string {
	tabs := []struct {
		label  string
		active bool
	}{
		{"Queue", m.showQueue},
		{"Search", m.isSearching},
		{"Now playing", !m.showQueue && !m.isSearching},
	}

	parts := make([]string, 0, len(tabs))
	for _, tab := range tabs {
		style := TabInactiveStyle
		if tab.active {
			style = TabActiveStyle
		}
		parts = append(parts, style.Render(tab.label))
	}

	return TabBarStyle.Width(width).Render(strings.Join(parts, "  "))
}
