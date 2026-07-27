package ui

import "fmt"

func renderSearch(m Model, width int) string {
	if m.isSearching {
		return SearchIconStyle.Render(IconSearch) + " " + SearchInputFocusedStyle.Render(m.searchInput.View())
	}

	if len(m.searchResults) > 0 {
		return SearchIconStyle.Render(IconSearch) + " " + TextDimStyle.Render(fmt.Sprintf("%d results | press / to search", len(m.searchResults)))
	}

	return SearchIconStyle.Render(IconSearch) + " " + SearchPlaceholderStyle.Render("search for a song...")
}
