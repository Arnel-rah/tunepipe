package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func renderPlayer(m Model, width int) string {
	if m.currentTrack == nil {
		return TextMutedStyle.Render(IconMusic + "  nothing playing")
	}

	var status string
	var statusStyle lipgloss.Style
	var icon string

	if m.isLoading {
		status = "LOADING"
		statusStyle = StatusLoadingStyle
		icon = IconLoading
	} else if m.isPlaying {
		status = "PLAYING"
		statusStyle = StatusPlayingStyle
		icon = IconPlay
	} else if m.elapsedTime.Seconds() >= m.currentTrack.Duration && m.currentTrack.Duration > 0 {
		status = "FINISHED"
		statusStyle = StatusIdleStyle
		icon = IconStop
	} else {
		status = "PAUSED"
		statusStyle = StatusPausedStyle
		icon = IconPause
	}

	maxLen := width - 4
	if maxLen < 10 {
		maxLen = 10
	}

	eyebrow := statusStyle.Render(icon + " " + status)
	title := NowPlayingTitleStyle.Render(truncate(m.currentTrack.Title, maxLen))
	artist := NowPlayingArtistStyle.Render(truncate(m.currentTrack.Uploader, maxLen))

	rail := RailStyle.Render(IconRail)
	textBlock := lipgloss.JoinVertical(lipgloss.Left, eyebrow, title, artist)

	return lipgloss.JoinHorizontal(lipgloss.Top, rail+" ", textBlock)
}
