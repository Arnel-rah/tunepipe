package ui

import (
	"strings"

	"tunepipe/internal/ytdlp"

	"github.com/charmbracelet/lipgloss"
)

func renderResults(m Model, width int) string {
	if len(m.searchResults) == 0 {
		return ""
	}

	maxItems := 8
	if m.height < 20 {
		maxItems = 5
	} else if m.height < 25 {
		maxItems = 6
	} else if m.height < 30 {
		maxItems = 8
	}

	lines := make([]string, 0, maxItems)

	for i := 0; i < len(m.searchResults) && i < maxItems; i++ {
		lines = append(lines,
			renderResultLine(
				m.searchResults[i],
				i == m.cursor,
				width,
			),
		)
	}

	return strings.Join(lines, "\n")
}

func renderResultLine(track ytdlp.Track, selected bool, width int) string {
	const (
		railWidth   = 2
		gap         = 2
		artistWidth = 20
	)

	lineStyle := RowStyle
	titleStyle := RowTitleStyle
	artistStyle := RowArtistStyle
	rail := "  "
	if selected {
		lineStyle = RowSelectedStyle
		titleStyle = RowTitleSelectedStyle
		artistStyle = RowArtistSelectedStyle
		rail = RailStyle.Render(IconRail) + " "
	}

	innerWidth := width - 2
	titleWidth := innerWidth - railWidth - gap - artistWidth
	if titleWidth < 5 {
		titleWidth = 5
	}
	if width < 50 {
		titleWidth = width - 25
		if titleWidth < 5 {
			titleWidth = 5
		}
	}

	title := truncate(track.Title, titleWidth)
	artist := truncate(track.Uploader, artistWidth)

	left := titleStyle.
		Width(titleWidth).
		MaxWidth(titleWidth).
		Render(title)

	right := artistStyle.
		Width(artistWidth).
		Align(lipgloss.Right).
		Render(artist)

	line := lipgloss.JoinHorizontal(
		lipgloss.Left,
		rail,
		left,
		"  ",
		right,
	)

	return lineStyle.Width(width).Render(line)
}
