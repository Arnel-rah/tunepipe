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
			Bold(true).
			Render(IconPlay)

	PauseIconStyle = lipgloss.NewStyle().
			Foreground(Amber).
			Bold(true).
			Render(IconPause)

	StopIconStyle = lipgloss.NewStyle().
			Foreground(TextMutedColor).
			Render(IconStop)

	LoadingIconStyle = lipgloss.NewStyle().
				Foreground(Violet).
				Bold(true).
				Render(IconLoading)

	MusicIconStyle = lipgloss.NewStyle().
			Foreground(Violet).
			Render(IconMusic)

	SearchIconStyle = lipgloss.NewStyle().
			Foreground(Violet).
			Render(IconSearch)

	ArrowIconStyle = lipgloss.NewStyle().
			Foreground(Violet).
			Render(IconArrow)

	BulletIconStyle = lipgloss.NewStyle().
			Foreground(TextMutedColor).
			Render(IconBullet)

	DiamondIconStyle = lipgloss.NewStyle().
				Foreground(Violet).
				Render(IconDiamond)

	CrossIconStyle = lipgloss.NewStyle().
			Foreground(Red).
			Render(IconCross)

	HeartIconStyle = lipgloss.NewStyle().
			Foreground(Red).
			Render(IconHeart)

	StarIconStyle = lipgloss.NewStyle().
			Foreground(Amber).
			Render(IconStar)
)
