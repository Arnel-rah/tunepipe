package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderFooter(m Model, width int) string {
	status := m.statusMsg
	if status == "" {
		status = "ready"
	}

	var style lipgloss.Style
	var icon string

	if m.isLoading {
		style = StatusLoadingStyle
		icon = IconLoading
		status = "loading..."
	} else if m.isPlaying {
		style = StatusPlayingStyle
		icon = IconPlay
		status = "playing"
	} else if m.statusKind == "error" {
		style = StatusErrorStyle
		icon = IconCross
	} else if m.currentTrack != nil && m.elapsedTime.Seconds() >= m.currentTrack.Duration {
		style = StatusIdleStyle
		icon = IconStop
		status = "finished"
	} else if len(m.searchResults) == 0 && m.currentTrack == nil {
		style = StatusIdleStyle
		icon = IconMusic
		status = "welcome"
	} else {
		style = StatusIdleStyle
		icon = IconBullet
		if status == "" {
			status = "ready"
		}
	}

	statusText := style.Render(icon + " " + strings.ToUpper(status))

	if len(m.queue) > 0 {
		queueText := AccentStyle.Render(fmt.Sprintf("%s %d queued", IconQueued, len(m.queue)))
		statusText = statusText + "  " + queueText
	}

	if width < 50 {
		return statusText
	}

	keys := []struct {
		key   string
		label string
	}{
		{"SPACE", "Play"},
		{"/", "Search"},
		{"↑↓", "Move"},
		{"A", "Toggle queue"},
		{"Q", "Queue"},
	}
	if m.showQueue {
		keys[3].label = "Clear queue"
		keys[4].label = "Results"
	}

	if width >= 80 {
		if m.showQueue {
			keys = []struct {
				key   string
				label string
			}{
				{"SPACE", "Play/Pause"},
				{"↑↓", "Navigate"},
				{"x", "Remove item"},
				{"ENTER", "Play now"},
				{"A", "Clear queue"},
				{"Q", "Results"},
			}
		} else {
			keys = []struct {
				key   string
				label string
			}{
				{"SPACE", "Play/Pause"},
				{"/", "Search"},
				{"↑↓", "Navigate"},
				{"A", "Toggle queue"},
				{"ENTER", "Select"},
				{"Q", "Queue"},
			}
		}
	}

	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s %s",
			KeyHintKeyStyle.Render(k.key),
			KeyHintTextStyle.Render(k.label),
		))
	}

	keysText := strings.Join(parts, " | ")

	if width < 70 && len(parts) > 3 {
		keysText = strings.Join(parts[:3], " | ")
	}

	gap := width - lipgloss.Width(statusText) - lipgloss.Width(keysText) - 2
	if gap < 1 {
		gap = 1
	}

	return fmt.Sprintf("%s%s%s", statusText, strings.Repeat(" ", gap), keysText)
}
