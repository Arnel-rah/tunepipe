package ui

import "github.com/charmbracelet/lipgloss"

var (
	BgColor     = lipgloss.Color("#0D0716")
	SurfaceColor = lipgloss.Color("#150C24")

	Violet     = lipgloss.Color("#8B5CF6")
	VioletDim  = lipgloss.Color("#5B4487")
	Magenta    = lipgloss.Color("#EC4899")
	White      = lipgloss.Color("#F5F2FF")
	Lavender   = lipgloss.Color("#A99FC4")
	LavenderMuted = lipgloss.Color("#544A6E")
	Red        = lipgloss.Color("#F5556C")
	Amber      = lipgloss.Color("#F7B955")

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

