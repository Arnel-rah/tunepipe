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
	} else {
		style = StatusIdleStyle
		icon = IconBullet
		if status == "" {
			status = "ready"
		}
	}

	statusText := style.Render(icon + " " + strings.ToUpper(status))

	keys := []struct {
		key   string
		label string
	}{
		{"SPACE", "Play/Pause"},
		{"/", "Search"},
		{"↑↓", "Navigate"},
		{"ENTER", "Select"},
		{"Q", "Quit"},
	}

	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s %s",
			KeyHintKeyStyle.Render(k.key),
			KeyHintTextStyle.Render(k.label),
		))
	}

	keysText := strings.Join(parts, "  ")

	gap := width - lipgloss.Width(statusText) - lipgloss.Width(keysText) - 2
	if gap < 1 {
		gap = 1
	}

	return fmt.Sprintf("%s%s%s", statusText, strings.Repeat(" ", gap), keysText)
}
