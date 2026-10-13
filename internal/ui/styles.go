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

	LogoStyle   = lipgloss.NewStyle().Foreground(Violet).Bold(true)
	AccentStyle = lipgloss.NewStyle().Foreground(PrimaryPink).Bold(true)

	StatLabelStyle = lipgloss.NewStyle().Foreground(TextMutedColor)
	StatValueStyle = lipgloss.NewStyle().Foreground(White).Bold(true)

	StatusPlayingStyle = lipgloss.NewStyle().Foreground(PrimaryPink).Bold(true)
	StatusPausedStyle  = lipgloss.NewStyle().Foreground(Amber).Bold(true)
	StatusErrorStyle   = lipgloss.NewStyle().Foreground(Red).Bold(true)
	StatusIdleStyle    = lipgloss.NewStyle().Foreground(TextMutedColor)
	StatusLoadingStyle = lipgloss.NewStyle().Foreground(Violet).Bold(true)

	ProgressBarStyle = lipgloss.NewStyle().Foreground(PrimaryPink)

	ProgressHeadStyle = lipgloss.NewStyle().Foreground(PrimaryPink).Bold(true)

	ProgressTrackStyle = lipgloss.NewStyle().Foreground(VioletDim)

	SearchInputStyle = lipgloss.NewStyle().
				Foreground(TextColor)

	SearchInputFocusedStyle = lipgloss.NewStyle().
				Foreground(White)

	SearchPlaceholderStyle = lipgloss.NewStyle().
				Foreground(TextMutedColor)

	RowStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(TextDimColor)

	RowSelectedStyle = lipgloss.NewStyle().
				Padding(0, 1).
				BorderLeft(true).
				BorderForeground(PrimaryPink).
				Foreground(White)

	RowTitleStyle          = lipgloss.NewStyle().Foreground(TextDimColor)
	RowTitleSelectedStyle  = lipgloss.NewStyle().Foreground(White).Bold(true)
	RowArtistStyle         = lipgloss.NewStyle().Foreground(TextMutedColor)
	RowArtistSelectedStyle = lipgloss.NewStyle().Foreground(PrimaryPink)

	KeyHintKeyStyle = lipgloss.NewStyle().
			Foreground(Violet).
			Bold(true)

	KeyHintTextStyle = lipgloss.NewStyle().
				Foreground(TextMutedColor)

	DividerStyle = lipgloss.NewStyle().Foreground(LavenderMuted)

	MutedStyle = lipgloss.NewStyle().Foreground(TextDimColor)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(White).
			BorderBottom(true).
			BorderForeground(VioletDim).
			Padding(0, 1)

	NowPlayingTitleStyle  = lipgloss.NewStyle().Foreground(White).Bold(true)
	NowPlayingArtistStyle = lipgloss.NewStyle().Foreground(TextDimColor)
	NowPlayingMetaStyle   = lipgloss.NewStyle().Foreground(LavenderMuted)

	RailStyle    = lipgloss.NewStyle().Foreground(PrimaryPink).Bold(true)
	RailDimStyle = lipgloss.NewStyle().Foreground(LavenderMuted)

	WelcomeTitleStyle    = LogoStyle
	WelcomeSubtitleStyle = SubtitleStyle
	WelcomeHintStyle     = lipgloss.NewStyle().Foreground(Amber).Bold(true)

	TopBarStyle = lipgloss.NewStyle().
			BorderLeft(true).
			BorderForeground(PrimaryPink).
			Padding(0, 1)

	TabBarStyle = lipgloss.NewStyle().
			Foreground(TextMutedColor).
			Padding(0, 0)

	TabActiveStyle = lipgloss.NewStyle().
			Foreground(Violet).
			Bold(true).
			Underline(true)

	TabInactiveStyle = lipgloss.NewStyle().
				Foreground(TextMutedColor).
				Padding(0, 0)

	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(TextMutedColor).
				Bold(true)

	ArtworkPanelStyle = lipgloss.NewStyle().
				BorderLeft(true).
				BorderForeground(VioletDim).
				Padding(0, 1)

	HeaderCardStyle = lipgloss.NewStyle().
			BorderLeft(true).
			BorderForeground(Violet).
			Padding(0, 1).
			MarginBottom(1)

	PlayerCardStyle = lipgloss.NewStyle().
			BorderLeft(true).
			BorderForeground(PrimaryPink).
			Padding(0, 1).
			MarginBottom(1)

	SearchCardStyle = lipgloss.NewStyle().
			BorderLeft(true).
			BorderForeground(Magenta).
			Padding(0, 1).
			MarginBottom(1)

	ResultsCardStyle = lipgloss.NewStyle().
				BorderLeft(true).
				BorderForeground(PrimaryPink).
				Padding(0, 1).
				MarginBottom(1)

	QueueCardStyle = lipgloss.NewStyle().
			BorderLeft(true).
			BorderForeground(VioletDim).
			Padding(0, 1).
			MarginBottom(1)

	ProgressCardStyle = lipgloss.NewStyle().
				BorderLeft(true).
				BorderForeground(LavenderMuted).
				Padding(0, 1).
				MarginBottom(1)

	FooterCardStyle = lipgloss.NewStyle().
			BorderTop(true).
			BorderForeground(VioletDim).
			Padding(1, 0)

	SectionLabelStyle = lipgloss.NewStyle().
				Foreground(TextMutedColor).
				Bold(true).
				Transform(strings.ToUpper)

	SearchShellStyle = lipgloss.NewStyle().
				BorderBottom(true).
				BorderForeground(Violet).
				Padding(0, 0)

	ResultsShellStyle = lipgloss.NewStyle().
				BorderTop(true).
				BorderBottom(true).
				BorderForeground(VioletDim).
				Padding(1, 0)
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
