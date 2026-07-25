package ui

import "github.com/charmbracelet/lipgloss"

// Palette inspired by rocksky.app — deep violet-black canvas,
// a single vivid violet accent, and a warm magenta for secondary emphasis.
// Flat by design: no card backgrounds, no boxed surfaces.
var (
	BgColor     = lipgloss.Color("#0D0716") // canvas
	SurfaceColor = lipgloss.Color("#150C24") // barely-there elevation (hover/input only)

	Violet     = lipgloss.Color("#8B5CF6") // primary accent
	VioletDim  = lipgloss.Color("#5B4487") // muted accent for rails/markers
	Magenta    = lipgloss.Color("#EC4899") // secondary accent (paused/heart)
	White      = lipgloss.Color("#F5F2FF")
	Lavender   = lipgloss.Color("#A99FC4") // dim text, tinted toward the palette
	LavenderMuted = lipgloss.Color("#544A6E") // faint text / rails / dividers
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

var RockskyTheme = Theme{
	Bg:        BgColor,
	Surface:   SurfaceColor,
	Violet:    Violet,
	VioletDim: VioletDim,
	Magenta:   Magenta,
	White:     White,
	Lavender:  Lavender,
	Muted:     LavenderMuted,
	Red:       Red,
	Amber:     Amber,
	Text:      TextColor,
	TextDim:   TextDimColor,
	TextMuted: TextMutedColor,
}
