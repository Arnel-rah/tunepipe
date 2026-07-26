package ui

import (
	"fmt"
	"strings"
)

func renderProgress(m Model, width int) string {
	if m.currentTrack == nil || m.currentTrack.Duration <= 0 {
		return ""
	}

	barWidth := width - 18
	if barWidth < 10 {
		barWidth = 10
	}

	var pct float64
	if m.elapsedTime.Seconds() >= m.currentTrack.Duration {
		pct = 1.0
	} else if m.currentTrack.Duration > 0 {
		pct = m.elapsedTime.Seconds() / m.currentTrack.Duration
		if pct > 1 {
			pct = 1
		}
		if pct < 0 {
			pct = 0
		}
	}

	filled := int(pct * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	if filled < 0 {
		filled = 0
	}

	elapsed := fmt.Sprintf("%02d:%02d", int(m.elapsedTime.Minutes()), int(m.elapsedTime.Seconds())%60)
	total := fmt.Sprintf("%02d:%02d", int(m.currentTrack.Duration/60), int(m.currentTrack.Duration)%60)

	var bar string
	if filled <= 0 {
		bar = ProgressTrackStyle.Render(strings.Repeat(CharTrack, barWidth))
	} else if filled >= barWidth {
		bar = ProgressBarStyle.Render(strings.Repeat(CharProgress, barWidth))
	} else {
		bar = ProgressBarStyle.Render(strings.Repeat(CharProgress, filled)) +
			ProgressTrackStyle.Render(strings.Repeat(CharTrack, barWidth-filled))
	}

	return fmt.Sprintf("%s %s %s",
		TextDimStyle.Render(elapsed),
		bar,
		TextDimStyle.Render(total),
	)
}
