package ui

import "fmt"

func renderSearch(m Model, width int) string {
	label := SectionLabelStyle.Render("01  SEARCH")
	if m.showQueue {
		if len(m.queue) == 0 {
			return SearchShellStyle.Width(width).Render(label + "  " + TextDimStyle.Render("queue empty · press Q to return"))
		}
		return SearchShellStyle.Width(width).Render(label + "  " + TextDimStyle.Render(fmt.Sprintf("%d queued · press Q to return", len(m.queue))))
	}

	if m.isSearching {
		return SearchShellStyle.Width(width).Render(label + "  " + SearchIconStyle.Render(IconSearch) + " " + SearchInputFocusedStyle.Render(m.searchInput.View()))
	}

	if len(m.searchResults) > 0 {
		return SearchShellStyle.Width(width).Render(label + "  " + TextDimStyle.Render(fmt.Sprintf("%d results · press / to search", len(m.searchResults))))
	}

	return SearchShellStyle.Width(width).Render(label + "  " + SearchPlaceholderStyle.Render("press / to search for a song..."))
}
