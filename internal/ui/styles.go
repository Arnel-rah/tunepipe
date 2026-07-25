package ui

import "github.com/charmbracelet/lipgloss"

var (
	TextPrimary   = lipgloss.Color("#F0F0F8")
	TextSecondary = lipgloss.Color("#B0B0C8")
	TextMuted     = lipgloss.Color("#707090")
	TextDim       = lipgloss.Color("#404060")

	AccentSky     = lipgloss.Color("#6C63FF")
	AccentIndigo  = lipgloss.Color("#818CF8")
	AccentEmerald = lipgloss.Color("#10B981")
	AccentAmber   = lipgloss.Color("#F59E0B")
	AccentRose    = lipgloss.Color("#EF4444")
)

var (
	TextStyle     = lipgloss.NewStyle().Foreground(TextPrimary)
	TextStyleDim  = lipgloss.NewStyle().Foreground(TextSecondary)
	MutedStyle    = lipgloss.NewStyle().Foreground(TextMuted)
	MutedStyleDim = lipgloss.NewStyle().Foreground(TextDim)

	TitleStyle      = lipgloss.NewStyle().Foreground(TextPrimary).Bold(true)
	LogoStyle       = lipgloss.NewStyle().Foreground(AccentSky).Bold(true)
	AccentTextStyle = lipgloss.NewStyle().Foreground(AccentSky).Bold(true)

	StatLabelStyle = lipgloss.NewStyle().Foreground(TextMuted)
	StatValueStyle = lipgloss.NewStyle().Foreground(AccentSky).Bold(true)

	DividerStyle = lipgloss.NewStyle().Foreground(TextDim)
)

var (
	SearchInputStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Padding(0, 1)

	SearchInputFocusedStyle = SearchInputStyle.
				Foreground(TextPrimary)
)

var (
	StatusTextPlaying = lipgloss.NewStyle().Foreground(AccentEmerald).Bold(true)
	StatusTextPaused  = lipgloss.NewStyle().Foreground(AccentAmber).Bold(true)
	StatusTextError   = lipgloss.NewStyle().Foreground(AccentRose).Bold(true)
	StatusTextIdle    = lipgloss.NewStyle().Foreground(TextMuted)

	ProgressBarStyle = lipgloss.NewStyle().Foreground(AccentSky)
)

var (
	RowStyle = lipgloss.NewStyle().
			Foreground(TextSecondary)

	RowSelectedStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Bold(true)

	KeyHintKeyStyle = lipgloss.NewStyle().
			Foreground(AccentSky).
			Bold(true)

	KeyHintTextStyle = lipgloss.NewStyle().
				Foreground(TextMuted)
)
