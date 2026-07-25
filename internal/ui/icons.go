package ui

import "github.com/charmbracelet/lipgloss"

const (
	IconPlay    = "▶"
	IconPause   = "⏸"
	IconStop    = "⏹"
	IconSearch  = "⌕"
	IconMusic   = "◆"
	IconNote    = "♪"
	IconRail    = "▎" // selection / now-playing rail marker, replaces arrow+card
	IconBullet  = "·"
	IconCheck   = "✓"
	IconCross   = "✕"
	IconLoading = "◐"
	IconDiamond = "◆"
	IconHeart   = "♥"
	IconShuffle = "⤨"
	IconRepeat  = "↻"
	IconVolume  = "◔"
	IconQueue   = "≡"

	CharProgress = "─"
	CharTrack    = "─"
	CharRule     = "─"
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
			Foreground(TextDimColor).
			Render(IconStop)

	SearchIconStyle = lipgloss.NewStyle().
			Foreground(TextMutedColor).
			Render(IconSearch)

	MusicIconStyle = lipgloss.NewStyle().
			Foreground(Violet).
			Render(IconMusic)

	RailIconStyle = lipgloss.NewStyle().
			Foreground(Violet).
			Bold(true).
			Render(IconRail)

	BulletIconStyle = lipgloss.NewStyle().
				Foreground(TextMutedColor).
				Render(IconBullet)

	CrossIconStyle = lipgloss.NewStyle().
			Foreground(Red).
			Render(IconCross)

	LoadingIconStyle = lipgloss.NewStyle().
				Foreground(Amber).
				Render(IconLoading)

	DiamondIconStyle = lipgloss.NewStyle().
				Foreground(Violet).
				Render(IconDiamond)

	HeartIconStyle = lipgloss.NewStyle().
			Foreground(Magenta).
			Render(IconHeart)

	ShuffleIconStyle = lipgloss.NewStyle().
				Foreground(Violet).
				Render(IconShuffle)

	RepeatIconStyle = lipgloss.NewStyle().
				Foreground(Violet).
				Render(IconRepeat)

	VolumeIconStyle = lipgloss.NewStyle().
				Foreground(TextMutedColor).
				Render(IconVolume)

	QueueIconStyle = lipgloss.NewStyle().
			Foreground(TextMutedColor).
			Render(IconQueue)
)
