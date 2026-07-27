package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	TextStyle      = lipgloss.NewStyle().Foreground(TextColor)
	TextDimStyle   = lipgloss.NewStyle().Foreground(TextDimColor)
	TextMutedStyle = lipgloss.NewStyle().Foreground(TextMutedColor)

	TitleStyle    = lipgloss.NewStyle().Foreground(White).Bold(true)
	SubtitleStyle = lipgloss.NewStyle().Foreground(TextDimColor)

	LogoStyle   = lipgloss.NewStyle().Foreground(PrimaryPink).Bold(true)
	AccentStyle = lipgloss.NewStyle().Foreground(PrimaryPink).Bold(true)

	StatLabelStyle = lipgloss.NewStyle().Foreground(TextMutedColor)
	StatValueStyle = lipgloss.NewStyle().Foreground(White).Bold(true)

	StatusPlayingStyle = lipgloss.NewStyle().Foreground(PrimaryPink).Bold(true)
	StatusPausedStyle  = lipgloss.NewStyle().Foreground(Amber).Bold(true)
	StatusErrorStyle   = lipgloss.NewStyle().Foreground(Red).Bold(true)
	StatusIdleStyle    = lipgloss.NewStyle().Foreground(TextMutedColor)
	StatusLoadingStyle = lipgloss.NewStyle().Foreground(Violet).Bold(true)

	ProgressBarStyle   = lipgloss.NewStyle().Foreground(PrimaryPink)
	ProgressTrackStyle = lipgloss.NewStyle().Foreground(LavenderMuted)

	SearchInputStyle = lipgloss.NewStyle().
				Foreground(TextColor)

	SearchInputFocusedStyle = lipgloss.NewStyle().
				Foreground(White)

	SearchPlaceholderStyle = lipgloss.NewStyle().
				Foreground(TextMutedColor)

	RowStyle = lipgloss.NewStyle().
			Padding(0, 1).
			BorderStyle(lipgloss.Border{Left: " "}).
			BorderLeft(true).
			BorderForeground(TextDimColor)

	RowSelectedStyle = lipgloss.NewStyle().
				Padding(0, 1).
				BorderStyle(lipgloss.Border{Left: "▍"}).
				BorderLeft(true).
				BorderForeground(PrimaryPink)

	RowTitleStyle          = lipgloss.NewStyle().Foreground(TextDimColor)
	RowTitleSelectedStyle  = lipgloss.NewStyle().Foreground(White).Bold(true)
	RowArtistStyle         = lipgloss.NewStyle().Foreground(TextMutedColor)
	RowArtistSelectedStyle = lipgloss.NewStyle().Foreground(PrimaryPink)

	KeyHintKeyStyle = lipgloss.NewStyle().
			Foreground(PrimaryPink).
			Bold(true)

	KeyHintTextStyle = lipgloss.NewStyle().
				Foreground(TextMutedColor)

	DividerStyle = lipgloss.NewStyle().Foreground(LavenderMuted)

	MutedStyle = lipgloss.NewStyle().Foreground(TextDimColor)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(White)

	NowPlayingTitleStyle  = lipgloss.NewStyle().Foreground(White).Bold(true)
	NowPlayingArtistStyle = lipgloss.NewStyle().Foreground(TextDimColor)

	RailStyle    = lipgloss.NewStyle().Foreground(PrimaryPink).Bold(true)
	RailDimStyle = lipgloss.NewStyle().Foreground(LavenderMuted)

	WelcomeTitleStyle    = LogoStyle
	WelcomeSubtitleStyle = SubtitleStyle
	WelcomeHintStyle     = lipgloss.NewStyle().Foreground(Amber).Bold(true)
)

func Divider(width int) string {
	if width < 1 {
		width = 1
	}
	return DividerStyle.Render(strings.Repeat(CharRule, width))
}

func DividerDot() string {
	return TextMutedStyle.Render(" " + IconBullet + " ")
}

func repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(s, n)
}
