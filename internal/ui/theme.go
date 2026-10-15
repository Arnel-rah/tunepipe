package ui

import "github.com/charmbracelet/lipgloss"

var (
	PrimaryPink   = lipgloss.Color("#f5c2e7") // pink
	Violet        = lipgloss.Color("#cba6f7") // mauve
	VioletDim     = lipgloss.Color("#585b70") // surface2
	Magenta       = lipgloss.Color("#f5c2e7") // pink
	White         = lipgloss.Color("#cdd6f4") // text
	Lavender      = lipgloss.Color("#a6adc8") // subtext0
	LavenderMuted = lipgloss.Color("#6c7086") // overlay0
	Red           = lipgloss.Color("#f38ba8") // red
	Amber         = lipgloss.Color("#f9e2af") // yellow

	Background = lipgloss.Color("#1e1e2e") // base
	Surface    = lipgloss.Color("#313244") // surface0

	TextColor      = White
	TextDimColor   = Lavender
	TextMutedColor = LavenderMuted
)

type Theme struct {
	Bg        lipgloss.Color
	Surface   lipgloss.Color
	Violet    lipgloss.Color
	VioletDim lipgloss.Color
	Magenta   lipgloss.Color
	White     lipgloss.Color
	Lavender  lipgloss.Color
	Muted     lipgloss.Color
	Red       lipgloss.Color
	Amber     lipgloss.Color
	Text      lipgloss.Color
	TextDim   lipgloss.Color
	TextMuted lipgloss.Color
}
