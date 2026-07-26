package ui

import "fmt"

func renderSearch(m Model, width int) string {
	var line string

	if m.isSearching {
		line = SearchIconStyle + " " + SearchInputFocusedStyle.Render(m.searchInput.View())
	} else if len(m.searchResults) > 0 {
		line = fmt.Sprintf("%s %s",
			SearchIconStyle,
			TextDimStyle.Render(fmt.Sprintf("%d results%s press / to search", len(m.searchResults), DividerDot())),
		)
	} else {
		line = SearchIconStyle + " " + SearchPlaceholderStyle.Render("search for a song...")
	}

	return line
}
