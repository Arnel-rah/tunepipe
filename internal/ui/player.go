package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func renderPlayer(m Model, width int) string {
	if m.currentTrack == nil {
		return PlayerCardStyle.Width(width).Render(TextMutedStyle.Render("NOW PLAYING  ·  " + IconMusic + "  ready"))
	}

	elapsed := fmt.Sprintf("%d:%02d", int(m.elapsedTime.Minutes()), int(m.elapsedTime.Seconds())%60)
	total := fmt.Sprintf("%d:%02d", int(m.currentTrack.Duration/60), int(m.currentTrack.Duration)%60)

	left := lipgloss.NewStyle().
		Width(16).
		Render(TextDimStyle.Render("NOW PLAYING\n" + elapsed + "/" + total))

	centerWidth := width - 32
	if centerWidth < 18 {
		centerWidth = 18
	}

	title := lipgloss.NewStyle().
		Width(centerWidth).
		Align(lipgloss.Center).
		Render(NowPlayingTitleStyle.Render(truncate(m.currentTrack.Title, centerWidth)))
	artist := lipgloss.NewStyle().
		Width(centerWidth).
		Align(lipgloss.Center).
		Render(NowPlayingArtistStyle.Render(truncate(m.currentTrack.Uploader, centerWidth)))
	center := lipgloss.JoinVertical(lipgloss.Center, title, artist)

	rightText := "QUEUE 0"
	if len(m.queue) > 0 {
		rightText = fmt.Sprintf("QUEUE %d", len(m.queue))
	}
	right := lipgloss.NewStyle().
		Width(14).
		Align(lipgloss.Right).
		Render(NowPlayingMetaStyle.Render(rightText))

	row := lipgloss.JoinHorizontal(lipgloss.Top, left, center, right)
	return PlayerCardStyle.Width(width).Render(row)
}
