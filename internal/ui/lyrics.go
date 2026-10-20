package ui

import (
	"strings"
)

func renderLyrics(m Model, width int) string {
	if !m.showLyrics || m.currentTrack == nil {
		return ""
	}
	if m.lyricsLoading {
		return LyricsPanelStyle.Width(width).Render(TextMutedStyle.Render("♪ loading lyrics..."))
	}
	if len(m.lyrics) == 0 {
		return LyricsPanelStyle.Width(width).Render(TextMutedStyle.Render("♪ lyrics unavailable"))
	}

	active := 0
	for i, line := range m.lyrics {
		if line.At <= m.elapsedTime {
			active = i
		} else {
			break
		}
	}
	start := active - 1
	if start < 0 {
		start = 0
	}
	end := start + 6
	if end > len(m.lyrics) {
		end = len(m.lyrics)
	}

	lines := []string{SectionLabelStyle.Render("LYRICS")}
	for i := start; i < end; i++ {
		style := TextDimStyle
		prefix := "  "
		if i == active {
			style = LyricsActiveStyle
			prefix = "› "
		}
		lines = append(lines, style.Render(prefix+truncate(m.lyrics[i].Text, width-4)))
	}
	return LyricsPanelStyle.Width(width).Render(strings.Join(lines, "\n"))
}
