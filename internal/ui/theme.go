package ui

import "github.com/charmbracelet/lipgloss"

var (
	PrimaryPink   = lipgloss.Color("#ff2876")
	Violet        = lipgloss.Color("#8d2dff")
	VioletDim     = lipgloss.Color("#5c2d91")
	Magenta       = lipgloss.Color("#ff6aa8")
	White         = lipgloss.Color("#f5efff")
	Lavender      = lipgloss.Color("#c8b5d0")
	LavenderMuted = lipgloss.Color("#7b6982")
	Red           = lipgloss.Color("#ff5b7a")
	Amber         = lipgloss.Color("#f7b955")

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
