package ui

import "fmt"

func renderSearch(m Model, width int) string {
	if m.showQueue {
		if len(m.queue) == 0 {
			return SearchIconStyle.Render(IconSearch) + " " + TextDimStyle.Render("queue empty | press Q to return")
		}
		return SearchIconStyle.Render(IconSearch) + " " + TextDimStyle.Render(fmt.Sprintf("%d queued | press Q to return", len(m.queue)))
	}

	if m.isSearching {
		return SearchIconStyle.Render(IconSearch) + " " + SearchInputFocusedStyle.Render(m.searchInput.View())
	}

	if len(m.searchResults) > 0 {
		return SearchIconStyle.Render(IconSearch) + " " + TextDimStyle.Render(fmt.Sprintf("%d results | press / to search", len(m.searchResults)))
	}

	return SearchIconStyle.Render(IconSearch) + " " + SearchPlaceholderStyle.Render("search for a song...")
}
