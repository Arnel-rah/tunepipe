package ui

import "github.com/charmbracelet/lipgloss"

var (
	PrimaryPink   = lipgloss.Color("#ff2bd6")
	Violet        = lipgloss.Color("#a855f7")
	VioletDim     = lipgloss.Color("#3b1f5e")
	Magenta       = lipgloss.Color("#ff00a8")
	Cyan          = lipgloss.Color("#00f5ff")
	White         = lipgloss.Color("#f4f7ff")
	Lavender      = lipgloss.Color("#b7b9d6")
	LavenderMuted = lipgloss.Color("#686b91")
	Red           = lipgloss.Color("#ff456e")
	Amber         = lipgloss.Color("#faff00")

	Background = lipgloss.Color("#0a0814")
	Surface    = lipgloss.Color("#151126")

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
	Cyan      lipgloss.Color
	White     lipgloss.Color
	Lavender  lipgloss.Color
	Muted     lipgloss.Color
	Red       lipgloss.Color
	Amber     lipgloss.Color
	Text      lipgloss.Color
	TextDim   lipgloss.Color
	TextMuted lipgloss.Color
}
