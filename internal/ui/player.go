package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func renderPlayer(m Model, width int) string {
	if m.currentTrack == nil {
		return TextMutedStyle.Render(IconMusic + "  ready to play")
	}

	var status string
	var statusStyle lipgloss.Style
	var icon string

	switch {
	case m.isLoading:
		status = "LOADING"
		statusStyle = StatusLoadingStyle
		icon = m.SpinnerFrame()
	case m.isPlaying:
		status = "PLAYING"
		statusStyle = StatusPlayingStyle
		icon = IconPlay
	case m.elapsedTime.Seconds() >= m.currentTrack.Duration && m.currentTrack.Duration > 0:
		status = "FINISHED"
		statusStyle = StatusIdleStyle
		icon = IconStop
	default:
		status = "PAUSED"
		statusStyle = StatusPausedStyle
		icon = IconPause
	}

	textWidth := width - 4
	if width < 50 {
		textWidth = width - 8
	}
	if textWidth < 10 {
		textWidth = 10
	}

	eyebrow := statusStyle.Render(icon + " " + status)
	title := NowPlayingTitleStyle.Render(truncate(m.currentTrack.Title, textWidth))
	artist := NowPlayingArtistStyle.Render(truncate(m.currentTrack.Uploader, textWidth))

	return lipgloss.JoinVertical(lipgloss.Left, eyebrow, title, artist)
}
