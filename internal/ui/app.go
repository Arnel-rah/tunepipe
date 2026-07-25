package ui

import (
	"fmt"

	"tunepipe/internal/player"
	"tunepipe/internal/ytdlp"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type searchResultsMsg struct {
	tracks []ytdlp.Track
	err    error
}

type trackReadyMsg struct {
	track     ytdlp.Track
	directURL string
	err       error
}

type Model struct {
	engine        *player.Engine
	searchInput   textinput.Model
	isSearching   bool
	searchResults []ytdlp.Track
	cursor        int
	currentTrack  *ytdlp.Track
	isPlaying     bool
	statusMsg     string
	width         int
	height        int
}

func NewModel(engine *player.Engine) Model {
	ti := textinput.New()
	ti.Placeholder = "Rechercher un morceau ou un artiste..."
	ti.CharLimit = 156
	ti.Width = 50
	ti.Prompt = "> "

	return Model{
		engine:      engine,
		searchInput: ti,
		statusMsg:   "Pret. Appuie sur '/' pour chercher.",
		width:       80,
		height:      24,
		isPlaying:   false,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func performSearch(query string) tea.Cmd {
	return func() tea.Msg {
		results, err := ytdlp.Search(query, 10)
		return searchResultsMsg{tracks: results, err: err}
	}
}

func fetchTrack(engine *player.Engine, track ytdlp.Track) tea.Cmd {
	return func() tea.Msg {
		directURL, err := ytdlp.FetchDirectURL(track.ID)
		if err == nil {
			err = engine.PlayURL(directURL)
			if err == nil {
				return trackReadyMsg{track: track, directURL: directURL, err: nil}
			}
			return trackReadyMsg{track: track, directURL: "", err: err}
		}
		return trackReadyMsg{track: track, directURL: "", err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var _ tea.Cmd

	if m.isSearching {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				query := m.searchInput.Value()
				if query != "" {
					m.isSearching = false
					m.statusMsg = "Recherche en cours avec yt-dlp..."
					m.searchInput.Blur()
					return m, performSearch(query)
				}
				m.isSearching = false
				m.searchInput.Blur()
				m.statusMsg = "Recherche annulee."
				return m, nil

			case "esc":
				m.isSearching = false
				m.searchInput.Blur()
				m.statusMsg = "Recherche annulee."
				return m, nil
			}
		}

		var updateCmd tea.Cmd
		m.searchInput, updateCmd = m.searchInput.Update(msg)
		return m, updateCmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case searchResultsMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Erreur recherche: %v", msg.err)
		} else {
			m.searchResults = msg.tracks
			m.cursor = 0
			m.statusMsg = fmt.Sprintf("%d resultats trouves.", len(msg.tracks))
		}
		return m, nil

	case trackReadyMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Erreur lecture: %v", msg.err)
			m.isPlaying = false
		} else {
			m.isPlaying = true
			m.statusMsg = fmt.Sprintf("Lecture en cours: %s", msg.track.Title)
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.engine != nil {
				m.engine.Stop()
			}
			return m, tea.Quit

		case "/":
			m.isSearching = true
			m.searchInput.Focus()
			m.searchInput.SetValue("")
			m.statusMsg = "Recherche: tapez votre requete puis Entree"
			return m, nil

		case "j", "down":
			if m.cursor < len(m.searchResults)-1 {
				m.cursor++
			}

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case " ":
			if m.engine != nil && m.currentTrack != nil {
				err := m.engine.TogglePause()
				if err != nil {
					m.statusMsg = fmt.Sprintf("Erreur pause: %v", err)
				} else {
					m.isPlaying = !m.isPlaying
					if m.isPlaying {
						m.statusMsg = "Lecture"
					} else {
						m.statusMsg = "Pause"
					}
				}
			} else {
				m.statusMsg = "Aucun morceau en cours de lecture"
			}

		case "enter":
			if len(m.searchResults) > 0 && m.cursor < len(m.searchResults) {
				selected := m.searchResults[m.cursor]
				m.currentTrack = &selected
				m.statusMsg = fmt.Sprintf("Extraction audio en cours: %s...", selected.Title)
				return m, fetchTrack(m.engine, selected)
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	mainWidth := m.width - 35
	if mainWidth < 20 {
		mainWidth = 20
	}
	playerWidth := m.width - 4
	if playerWidth < 20 {
		playerWidth = 20
	}

	var searchContent string
	if m.isSearching {
		searchBox := SearchInputStyle.Width(mainWidth - 4).Render(m.searchInput.View())
		searchContent = SearchModeStyle + "\n\n" + searchBox
		searchContent += "\n\n" + MutedStyle.Render("Appuyez sur Entree pour valider, Echap pour annuler")
	} else if len(m.searchResults) == 0 {
		searchContent = MutedStyle.Render("Appuie sur '/' pour lancer une recherche.")
	} else {
		for i, track := range m.searchResults {
			cursorStr := " "
			item := fmt.Sprintf("%d. %s [%s]", i+1, track.Title, track.Uploader)
			if m.cursor == i {
				cursorStr = ">"
				itemStr := fmt.Sprintf("%s %s", cursorStr, item)
				searchContent += SelectedItemStyle.Render(itemStr) + "\n"
			} else {
				itemStr := fmt.Sprintf("%s %s", cursorStr, item)
				searchContent += itemStr + "\n"
			}
		}
	}

	mainPane := MainPaneStyle.
		Width(mainWidth).
		Height(12).
		Render("RECHERCHE YOUTUBE\n\n" + searchContent)

	sidePane := SidebarStyle.
		Width(30).
		Height(12).
		Render("CONTROLES\n\n[/] Chercher\n[j/k] Naviguer\n[Enter] Lire\n[Space] Pause\n[q] Quitter")

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, mainPane, sidePane)

	nowPlaying := "Aucun morceau en lecture"
	if m.currentTrack != nil {
		nowPlaying = fmt.Sprintf(" %s - %s", m.currentTrack.Title, m.currentTrack.Uploader)
	}

	statusStyle := StatusPausedStyle
	if m.isPlaying {
		statusStyle = StatusPlayingStyle
	}

	playerBar := PlayerBarStyle.
		Width(playerWidth).
		Render(fmt.Sprintf("%s\n%s", statusStyle.Render(nowPlaying), MutedStyle.Render(m.statusMsg)))

	return lipgloss.JoinVertical(lipgloss.Left, topRow, playerBar)
}
