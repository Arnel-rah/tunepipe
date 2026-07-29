package ui

import (
	"fmt"
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

	start := 0
	if len(m.searchResults) > maxItems && m.cursor >= maxItems {
		start = m.cursor - maxItems + 1
		if start > len(m.searchResults)-maxItems {
			start = len(m.searchResults) - maxItems
		}
	}
	end := start + maxItems
	if end > len(m.searchResults) {
		end = len(m.searchResults)
	}

	lines := make([]string, 0, maxItems+1)
	lines = append(lines, renderTableHeader(width))

	for i := start; i < end; i++ {
		track := m.searchResults[i]
		lines = append(lines,
			renderResultLine(track, i == m.cursor, width),
		)
	}

	return strings.Join(lines, "\n")
}

func renderTableHeader(width int) string {
	titleWidth, artistWidth, durationWidth := tableColumnWidths(width)
	return TableHeaderStyle.Render(fmt.Sprintf("%-*s  %-*s  %*s", titleWidth, "Title", artistWidth, "Artist", durationWidth, "Duration"))
}

func renderResultLine(track ytdlp.Track, selected bool, width int) string {
	titleWidth, artistWidth, durationWidth := tableColumnWidths(width)

	lineStyle := RowStyle
	titleStyle := RowTitleStyle
	artistStyle := RowArtistStyle
	if selected {
		lineStyle = RowSelectedStyle
		titleStyle = RowTitleSelectedStyle
		artistStyle = RowArtistSelectedStyle
	}

	title := truncate(track.Title, titleWidth)
	artist := truncate(track.Uploader, artistWidth)
	duration := fmt.Sprintf("%d:%02d", int(track.Duration/60), int(track.Duration)%60)

	left := titleStyle.
		Width(titleWidth).
		MaxWidth(titleWidth).
		Render(title)

	right := artistStyle.
		Width(artistWidth).
		Align(lipgloss.Right).
		Render(artist)

	durationText := lipgloss.NewStyle().
		Width(durationWidth).
		Align(lipgloss.Right).
		Render(TextMutedStyle.Render(duration))

	line := lipgloss.JoinHorizontal(lipgloss.Left, left, "  ", right, "  ", durationText)

	return lineStyle.Width(width).Render(line)
}

func tableColumnWidths(width int) (int, int, int) {
	durationWidth := 7
	artistWidth := 24
	if width < 80 {
		artistWidth = 18
	}
	if width < 65 {
		artistWidth = 14
	}
	titleWidth := width - artistWidth - durationWidth - 6
	if titleWidth < 12 {
		titleWidth = 12
	}
	return titleWidth, artistWidth, durationWidth
}
