package ui

import "strings"

func renderArtwork(width int) string {
	lines := []string{
		"      ╭────────────╮",
		"     ╱              ╲",
		"    │   ●      ●     │",
		"    │                │",
		"    │      ▣▣▣       │",
		"    │                │",
		"    │   ●        ●   │",
		"     ╲              ╱",
		"      ╰────────────╯",
	}

	return ArtworkPanelStyle.Width(width).Render(strings.Join(lines, "\n"))
}
