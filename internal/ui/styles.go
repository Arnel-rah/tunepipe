package ui

import "github.com/charmbracelet/lipgloss"

var (
	BackgroundColor = lipgloss.Color("#0A0A0F")
	SurfaceColor    = lipgloss.Color("#15151F")
	SurfaceLight    = lipgloss.Color("#232333")
	SurfaceHover    = lipgloss.Color("#2E2E42")

	PrimaryColor    = lipgloss.Color("#8B5CF6")
	PrimaryColorDim = lipgloss.Color("#5B3A9E")
	AccentColor     = lipgloss.Color("#00D4FF")
	SecondaryColor  = lipgloss.Color("#FF6B6B")
	MutedColor      = lipgloss.Color("#6B6B7B")
	MutedColorDim   = lipgloss.Color("#45454F")
	TextColor       = lipgloss.Color("#E4E4E8")
	TextColorBright = lipgloss.Color("#FFFFFF")

	NeonGreen  = lipgloss.Color("#00FF7F")
	NeonPink   = lipgloss.Color("#FF1493")
	NeonYellow = lipgloss.Color("#FFD700")
	GoldColor  = lipgloss.Color("#FBBF24")

	BorderRadiusStyle = lipgloss.RoundedBorder()

	MutedStyle      = lipgloss.NewStyle().Foreground(MutedColor)
	TextStyle       = lipgloss.NewStyle().Foreground(TextColor)
	TitleStyle      = lipgloss.NewStyle().Foreground(TextColorBright).Bold(true)
	SubtitleStyle   = lipgloss.NewStyle().Foreground(MutedColor).Italic(true)
	AccentTextStyle = lipgloss.NewStyle().Foreground(AccentColor).Bold(true)

	MainPaneStyle = lipgloss.NewStyle().
			Background(SurfaceColor).
			Border(BorderRadiusStyle).
			BorderForeground(PrimaryColorDim).
			Padding(1, 2)

	MainPaneFocusedStyle = MainPaneStyle.
				BorderForeground(PrimaryColor)

	SidebarStyle = lipgloss.NewStyle().
			Background(SurfaceColor).
			Border(BorderRadiusStyle).
			BorderForeground(MutedColorDim).
			Padding(1, 2)

	SidebarFocusedStyle = SidebarStyle.
				BorderForeground(AccentColor)

	PlayerBarStyle = lipgloss.NewStyle().
			Background(SurfaceColor).
			Border(BorderRadiusStyle).
			BorderForeground(AccentColor).
			Padding(1, 2)

	SearchInputStyle = lipgloss.NewStyle().
				Background(SurfaceLight).
				Foreground(TextColor).
				Padding(0, 2).
				Border(BorderRadiusStyle).
				BorderForeground(PrimaryColorDim)

	SearchInputFocusedStyle = SearchInputStyle.
				BorderForeground(PrimaryColor).
				Foreground(TextColorBright)

	SearchModeStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(TextColorBright).
				Background(PrimaryColor).
				Bold(true).
				Padding(0, 1)

	HoveredItemStyle = lipgloss.NewStyle().
				Foreground(TextColorBright).
				Background(SurfaceHover).
				Padding(0, 1)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Padding(0, 1)

	StatusPlayingStyle = lipgloss.NewStyle().
				Foreground(NeonGreen).
				Bold(true)

	StatusPausedStyle = lipgloss.NewStyle().
				Foreground(NeonYellow).
				Bold(true)

	StatusErrorStyle = lipgloss.NewStyle().
				Foreground(SecondaryColor).
				Bold(true)

	StatusLoadingStyle = lipgloss.NewStyle().
				Foreground(AccentColor).
				Bold(true)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true).
			Background(SurfaceColor).
			Padding(0, 2)

	CounterStyle = lipgloss.NewStyle().
			Foreground(GoldColor).
			Bold(true)

	TagStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Background(SurfaceLight).
			Padding(0, 1)

	BadgeStyle = lipgloss.NewStyle().
			Foreground(TextColorBright).
			Background(PrimaryColor).
			Bold(true).
			Padding(0, 1)

	ProgressBarStyle = lipgloss.NewStyle().
				Foreground(PrimaryColor)

	ProgressBgStyle = lipgloss.NewStyle().
			Foreground(SurfaceLight)

	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(MutedColor)
)

func Divider() string {
	return lipgloss.NewStyle().Foreground(MutedColor).Render(" • ")
}
