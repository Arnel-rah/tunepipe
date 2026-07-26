package ui

import "github.com/charmbracelet/lipgloss"

var (
	TextStyle      = lipgloss.NewStyle().Foreground(TextColor)
	TextDimStyle   = lipgloss.NewStyle().Foreground(TextDimColor)
	TextMutedStyle = lipgloss.NewStyle().Foreground(TextMutedColor)

	TitleStyle    = lipgloss.NewStyle().Foreground(White).Bold(true)
	SubtitleStyle = lipgloss.NewStyle().Foreground(TextDimColor)

	LogoStyle   = lipgloss.NewStyle().Foreground(Violet).Bold(true)
	AccentStyle = lipgloss.NewStyle().Foreground(Violet).Bold(true)

	StatLabelStyle = lipgloss.NewStyle().Foreground(TextMutedColor)
	StatValueStyle = lipgloss.NewStyle().Foreground(White).Bold(true)

	StatusPlayingStyle = lipgloss.NewStyle().Foreground(Violet).Bold(true)
	StatusPausedStyle  = lipgloss.NewStyle().Foreground(Amber).Bold(true)
	StatusErrorStyle   = lipgloss.NewStyle().Foreground(Red).Bold(true)
	StatusIdleStyle    = lipgloss.NewStyle().Foreground(TextMutedColor)
	StatusLoadingStyle = lipgloss.NewStyle().Foreground(Violet).Bold(true)

	ProgressBarStyle   = lipgloss.NewStyle().Foreground(Violet)
	ProgressTrackStyle = lipgloss.NewStyle().Foreground(LavenderMuted)

	SearchInputStyle = lipgloss.NewStyle().
				Foreground(TextColor)

	SearchInputFocusedStyle = lipgloss.NewStyle().
				Foreground(White)

	SearchPlaceholderStyle = lipgloss.NewStyle().
				Foreground(TextMutedColor)

	RowStyle = lipgloss.NewStyle().
			Padding(0, 1)

	RowSelectedStyle = lipgloss.NewStyle().
				Padding(0, 1)

	RowTitleStyle          = lipgloss.NewStyle().Foreground(TextDimColor)
	RowTitleSelectedStyle  = lipgloss.NewStyle().Foreground(White).Bold(true)
	RowArtistStyle         = lipgloss.NewStyle().Foreground(TextMutedColor)
	RowArtistSelectedStyle = lipgloss.NewStyle().Foreground(Violet)

	KeyHintKeyStyle = lipgloss.NewStyle().
			Foreground(Violet).
			Bold(true)

	KeyHintTextStyle = lipgloss.NewStyle().
				Foreground(TextMutedColor)

	DividerStyle = lipgloss.NewStyle().Foreground(LavenderMuted)

	MutedStyle = lipgloss.NewStyle().Foreground(TextDimColor)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(White)

	NowPlayingTitleStyle  = lipgloss.NewStyle().Foreground(White).Bold(true)
	NowPlayingArtistStyle = lipgloss.NewStyle().Foreground(TextDimColor)

	RailStyle    = lipgloss.NewStyle().Foreground(Violet).Bold(true)
	RailDimStyle = lipgloss.NewStyle().Foreground(LavenderMuted)
)

func Divider(width int) string {
	if width < 1 {
		width = 1
	}
	return DividerStyle.Render(repeat(CharRule, width))
}

func DividerDot() string {
	return TextMutedStyle.Render(" " + IconBullet + " ")
}

func GetStatusIcon(status string) string {
	switch status {
	case "playing":
		return PlayIconStyle
	case "paused":
		return PauseIconStyle
	case "loading":
		return LoadingIconStyle
	case "finished":
		return StopIconStyle
	case "error":
		return CrossIconStyle
	default:
		return BulletIconStyle
	}
}

func repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}
	out := make([]byte, 0, len(s)*n)
	for i := 0; i < n; i++ {
		out = append(out, s...)
	}
	return string(out)
}
