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
	switch m.statusKind {
	case "playing":
		style = StatusPlayingStyle
		icon = IconPlay
	case "paused":
		style = StatusPausedStyle
		icon = IconPause
	case "error":
		style = StatusErrorStyle
		icon = IconCross
	default:
		style = StatusIdleStyle
		icon = IconBullet
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

	keysText := strings.Join(parts, "   ")

	gap := width - lipgloss.Width(statusText) - lipgloss.Width(keysText)
	if gap < 1 {
		gap = 1
	}

	content := fmt.Sprintf("%s%s%s", statusText, strings.Repeat(" ", gap), keysText)
	return content
}
