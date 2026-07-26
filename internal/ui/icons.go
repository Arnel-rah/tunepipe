package ui

import "github.com/charmbracelet/lipgloss"

const (
	IconPlay    = "▶"
	IconPause   = "⏸"
	IconStop    = "⏹"
	IconSearch  = "⌕"
	IconMusic   = "◆"
	IconNote    = "♪"
	IconRail    = "▎"
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
			Background(BgColor).
			Render(IconBullet)

	CrossIconStyle = lipgloss.NewStyle().
			Foreground(Red).
			Background(BgColor).
			Render(IconCross)

	LoadingIconStyle = lipgloss.NewStyle().
				Foreground(Amber).
				Background(BgColor).
				Render(IconLoading)

	DiamondIconStyle = lipgloss.NewStyle().
				Foreground(Violet).
				Background(BgColor).
				Render(IconDiamond)

	HeartIconStyle = lipgloss.NewStyle().
			Foreground(Magenta).
			Background(BgColor).
			Render(IconHeart)

	ShuffleIconStyle = lipgloss.NewStyle().
				Foreground(Violet).
				Background(BgColor).
				Render(IconShuffle)

	RepeatIconStyle = lipgloss.NewStyle().
			Foreground(Violet).
			Background(BgColor).
			Render(IconRepeat)

	VolumeIconStyle = lipgloss.NewStyle().
			Foreground(TextMutedColor).
			Background(BgColor).
			Render(IconVolume)

	QueueIconStyle = lipgloss.NewStyle().
			Foreground(TextMutedColor).
			Background(BgColor).
			Render(IconQueue)
)
