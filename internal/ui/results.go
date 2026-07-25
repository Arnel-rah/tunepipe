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

	lines := make([]string, 0, 8)

	for i := 0; i < len(m.searchResults) && i < 8; i++ {
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
		artistWidth = 22
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

	innerWidth := width - lineStyle.GetHorizontalFrameSize()
	titleWidth := innerWidth - railWidth - gap - artistWidth
	if titleWidth < 10 {
		titleWidth = 10
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
