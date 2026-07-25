package ui

import "github.com/charmbracelet/lipgloss"

var (
	PrimaryColor   = lipgloss.Color("#8A2BE2")
	SecondaryColor = lipgloss.Color("#00FFFF")
	MutedColor     = lipgloss.Color("#555555")
	ActiveColor    = lipgloss.Color("#00FF7F")
	BgHeaderColor  = lipgloss.Color("#1A1A24")

	MutedStyle = lipgloss.NewStyle().Foreground(MutedColor)

	MainPaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(0, 1)

	SidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(MutedColor).
			Padding(0, 1)

	PlayerBarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SecondaryColor).
			Padding(0, 1)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(PrimaryColor).
				Bold(true)

	StatusPlayingStyle = lipgloss.NewStyle().Foreground(ActiveColor).Bold(true)
	StatusPausedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFCC00")).Bold(true)

	SearchInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#2A2A3A")).
				Padding(0, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#8A2BE2"))

	SearchModeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF")).
			Bold(true).
			Render("RECHERCHE")

	NormalModeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555")).
			Render("NORMAL")
)
