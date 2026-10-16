package ui

import (
	"fmt"
	"strings"
	"time"
)

func formatClock(d time.Duration, withHours bool) string {
	if d < 0 {
		d = 0
	}
	s := int(d.Seconds())
	h, m, sec := s/3600, (s%3600)/60, s%60
	if withHours {
		return fmt.Sprintf("%d:%02d:%02d", h, m, sec)
	}
	return fmt.Sprintf("%02d:%02d", m, sec)
}

func renderProgress(m Model, width int) string {
	if m.currentTrack == nil || m.currentTrack.Duration <= 0 {
		return ""
	}

	total := time.Duration(m.currentTrack.Duration * float64(time.Second))
	elapsed := m.elapsedTime
	if elapsed < 0 {
		elapsed = 0
	}
	if elapsed > total {
		elapsed = total
	}

	withHours := total >= time.Hour
	elapsedStr := formatClock(elapsed, withHours)
	totalStr := formatClock(total, withHours)

	showPercent := width >= 64
	percentStr := ""
	if showPercent {
		percentStr = fmt.Sprintf(" %3d%%", int(float64(elapsed)/float64(total)*100))
	}

	barWidth := width - len(elapsedStr) - len(totalStr) - len(percentStr) - 2
	if barWidth < 10 {
		barWidth = 10
	}

	pct := float64(elapsed) / float64(total)

	head := int(pct * float64(barWidth-1))
	if head < 0 {
		head = 0
	}
	if head > barWidth-1 {
		head = barWidth - 1
	}

	bar := ProgressBarStyle.Render(strings.Repeat(CharProgress, head)) +
		ProgressHeadStyle.Render(CharHead) +
		ProgressTrackStyle.Render(strings.Repeat(CharTrack, barWidth-head-1))

	result := fmt.Sprintf("%s %s %s",
		TextDimStyle.Render(elapsedStr),
		bar,
		TextDimStyle.Render(totalStr),
	)
	if showPercent {
		result += ProgressPercentStyle.Render(percentStr)
	}
	return result
}
