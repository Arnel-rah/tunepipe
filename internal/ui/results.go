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

	if m.searchPending {
		return TextDimStyle.Render(m.SpinnerFrame() + " searching...")
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
		track := m.searchResults[i]
		lines = append(lines,
			renderResultLine(
				track,
				i == m.cursor,
				m.IsQueued(track.ID),
				true,
				width,
			),
		)
	}

	return strings.Join(lines, "\n")
}

func renderResultLine(track ytdlp.Track, selected bool, queued bool, showQueuedBadge bool, width int) string {
	const (
		gap         = 2
		artistWidth = 20
	)

	lineStyle := RowStyle
	titleStyle := RowTitleStyle
	artistStyle := RowArtistStyle
	if selected {
		lineStyle = RowSelectedStyle
		titleStyle = RowTitleSelectedStyle
		artistStyle = RowArtistSelectedStyle
	}

	badge := ""
	badgeWidth := 0
	if queued && showQueuedBadge {
		badge = AccentStyle.Render(IconQueued) + " "
		badgeWidth = lipgloss.Width(badge)
	}

	innerWidth := width - 2
	titleWidth := innerWidth - gap - artistWidth - badgeWidth
	if titleWidth < 5 {
		titleWidth = 5
	}
	if width < 50 {
		titleWidth = width - 25 - badgeWidth
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
		badge,
		left,
		"  ",
		right,
	)

	return lineStyle.Width(width).Render(line)
}
