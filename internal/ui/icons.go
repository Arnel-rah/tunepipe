package ui

import "github.com/charmbracelet/lipgloss"

const (
	IconPlay    = "▶"
	IconPause   = "❚❚"
	IconStop    = "■"
	IconLoading = "◐"
	IconMusic   = "♫"
	IconNote    = "♪"
	IconRail    = "▎"
	IconSearch  = "⌕"
	IconArrow   = "▸"
	IconBullet  = "●"
	IconDiamond = "✦"
	IconCheck   = "✓"
	IconCross   = "✕"
	IconHeart   = "♥"
	IconStar    = "★"

	CharRule     = "─"
	CharProgress = "█"
	CharTrack    = "░"
)

var (
	PlayIconStyle = lipgloss.NewStyle().
			Foreground(Violet).
			Bold(true)

	PauseIconStyle = lipgloss.NewStyle().
			Foreground(Amber).
			Bold(true)

	StopIconStyle = lipgloss.NewStyle().
			Foreground(TextMutedColor)

	LoadingIconStyle = lipgloss.NewStyle().
				Foreground(Violet).
				Bold(true)

	MusicIconStyle = lipgloss.NewStyle().
			Foreground(Violet)

	SearchIconStyle = lipgloss.NewStyle().
			Foreground(Violet)

	ArrowIconStyle = lipgloss.NewStyle().
			Foreground(Violet)

	BulletIconStyle = lipgloss.NewStyle().
			Foreground(TextMutedColor)

	DiamondIconStyle = lipgloss.NewStyle().
				Foreground(Violet)

	CrossIconStyle = lipgloss.NewStyle().
			Foreground(Red)

	HeartIconStyle = lipgloss.NewStyle().
			Foreground(Red)

	StarIconStyle = lipgloss.NewStyle().
			Foreground(Amber)
)

func GetStatusIcon(status string) string {
	switch status {
	case "playing":
		return PlayIconStyle.Render(IconPlay)
	case "paused":
		return PauseIconStyle.Render(IconPause)
	case "loading":
		return LoadingIconStyle.Render(IconLoading)
	case "finished":
		return StopIconStyle.Render(IconStop)
	case "error":
		return CrossIconStyle.Render(IconCross)
	default:
		return BulletIconStyle.Render(IconBullet)
	}
}
