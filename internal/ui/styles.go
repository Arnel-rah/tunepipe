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

const (
	IconPlay      = "▶"
	IconPause     = "⏸"
	IconStop      = "⏹"
	IconNext      = "⏭"
	IconPrev      = "⏮"
	IconSearch    = "⌕"
	IconMusic     = "♫"
	IconNote      = "♪"
	IconStar      = "★"
	IconBullet    = "●"
	IconArrow     = "▸"
	IconDiamond   = "✦"
	IconCheck     = "✓"
	IconCross     = "✗"
	IconHeart     = "♥"
	IconSpeaker   = "♩"
	IconHeadphone = ""
	IconFolder    = ""
	IconFile      = ""
	IconGit       = ""
	IconTwitter   = ""
	IconGear      = ""
	IconHome      = ""
	IconUser      = ""
	IconClock     = ""
	IconCalendar  = ""
	IconList      = ""
	IconMenu      = ""
	IconSettings  = ""
)

// Styles avec icônes intégrées
var (
	PlayIconStyle = lipgloss.NewStyle().
			Foreground(AccentEmerald).
			Bold(true).
			Render(IconPlay)

	PauseIconStyle = lipgloss.NewStyle().
			Foreground(AccentAmber).
			Bold(true).
			Render(IconPause)

	StopIconStyle = lipgloss.NewStyle().
			Foreground(TextMuted).
			Render(IconStop)

	SearchIconStyle = lipgloss.NewStyle().
			Foreground(AccentSky).
			Render(IconSearch)

	MusicIconStyle = lipgloss.NewStyle().
			Foreground(AccentIndigo).
			Render(IconMusic)

	ArrowIconStyle = lipgloss.NewStyle().
			Foreground(TextDim).
			Render(IconArrow)

	DiamondIconStyle = lipgloss.NewStyle().
				Foreground(AccentSky).
				Render(IconDiamond)

	BulletIconStyle = lipgloss.NewStyle().
			Foreground(TextDim).
			Render(IconBullet)

	HeartIconStyle = lipgloss.NewStyle().
			Foreground(AccentRose).
			Render(IconHeart)

	SpeakerIconStyle = lipgloss.NewStyle().
				Foreground(AccentIndigo).
				Render(IconSpeaker)

	NoteIconStyle = lipgloss.NewStyle().
			Foreground(AccentSky).
			Render(IconNote)

	StarIconStyle = lipgloss.NewStyle().
			Foreground(AccentAmber).
			Render(IconStar)

	CheckIconStyle = lipgloss.NewStyle().
			Foreground(AccentEmerald).
			Render(IconCheck)

	CrossIconStyle = lipgloss.NewStyle().
			Foreground(AccentRose).
			Render(IconCross)
)

func GetPlayIcon(playing bool) string {
	if playing {
		return PlayIconStyle
	}
	return PauseIconStyle
}

func GetStatusIcon(state string) string {
	switch state {
	case "playing":
		return PlayIconStyle
	case "paused":
		return PauseIconStyle
	case "error":
		return CrossIconStyle
	default:
		return BulletIconStyle
	}
}

func GetSearchIcon(focused bool) string {
	if focused {
		return SearchIconStyle
	}
	return SearchIconStyle
}

func GetArrowIcon(selected bool) string {
	if selected {
		return ArrowIconStyle
	}
	return ArrowIconStyle
}

const (
	CharProgress = "━"
	CharTrack    = "─"
	CharDot      = "●"
)

const (
	BorderTop    = "─"
	BorderBottom = "─"
	BorderLeft   = "│"
	BorderRight  = "│"
	CornerTL     = "┌"
	CornerTR     = "┐"
	CornerBL     = "└"
	CornerBR     = "┘"
	LineH        = "─"
	LineV        = "│"
	LineCross    = "┼"
)

var (
	ConfirmStyle = lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(AccentSky).
		Foreground(TextPrimary)
)
