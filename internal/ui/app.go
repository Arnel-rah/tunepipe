package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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

type prefetchTickMsg struct {
	gen     int
	videoID string
}

type searchTickMsg struct {
	gen   int
	query string
}

type tickMsg time.Time

type Model struct {
	engine        *player.Engine
	cache         *ytdlp.URLCache
	searchInput   textinput.Model
	isSearching   bool
	searchResults []ytdlp.Track
	cursor        int
	selectionGen  int
	searchGen     int
	currentTrack  *ytdlp.Track
	isPlaying     bool
	statusMsg     string
	width         int
	height        int

	totalMinutes   int
	currentMinutes float64
	startTime      time.Time
	elapsedTime    time.Duration
	scrobbles      int
}

func NewModel(engine *player.Engine) Model {
	ti := textinput.New()
	ti.Placeholder = "Rechercher un morceau ou un artiste..."
	ti.CharLimit = 156
	ti.Width = 50
	ti.Prompt = "> "
	ti.TextStyle = lipgloss.NewStyle().Foreground(TextColor)

	return Model{
		engine:       engine,
		cache:        ytdlp.NewURLCache(5 * time.Minute),
		searchInput:  ti,
		statusMsg:    "Ready to play • Press '/' to search",
		width:        80,
		height:       24,
		isPlaying:    false,
		totalMinutes: 0,
		scrobbles:    0,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		tickCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func debouncePrefetch(gen int, videoID string) tea.Cmd {
	return tea.Tick(150*time.Millisecond, func(time.Time) tea.Msg {
		return prefetchTickMsg{gen: gen, videoID: videoID}
	})
}

func debounceSearch(gen int, query string) tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(time.Time) tea.Msg {
		return searchTickMsg{gen: gen, query: query}
	})
}

func performSearch(query string) tea.Cmd {
	return func() tea.Msg {
		results, err := ytdlp.Search(query, 10)
		return searchResultsMsg{tracks: results, err: err}
	}
}

func fetchTrack(engine *player.Engine, cache *ytdlp.URLCache, track ytdlp.Track) tea.Cmd {
	return func() tea.Msg {
		directURL, err := cache.Get(track.ID)
		if err != nil {
			return trackReadyMsg{track: track, directURL: "", err: err}
		}
		if err := engine.PlayURL(directURL); err != nil {
			return trackReadyMsg{track: track, directURL: "", err: err}
		}
		return trackReadyMsg{track: track, directURL: directURL, err: nil}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		if m.isPlaying && m.currentTrack != nil {
			m.elapsedTime += time.Second
			m.currentMinutes = m.elapsedTime.Minutes()

			if int(m.elapsedTime.Seconds())%10 == 0 {
				m.scrobbles++
			}
		}
		return m, tickCmd()

	case searchResultsMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Error: %v", msg.err)
			return m, nil
		}
		m.searchResults = msg.tracks
		m.cursor = 0
		m.statusMsg = fmt.Sprintf("%d results found", len(msg.tracks))

		prefetchCount := 3
		if len(msg.tracks) < prefetchCount {
			prefetchCount = len(msg.tracks)
		}
		for i := 0; i < prefetchCount; i++ {
			m.cache.Prefetch(msg.tracks[i].ID)
		}
		return m, nil

	case searchTickMsg:
		if msg.gen == m.searchGen && m.isSearching {
			return m, performSearch(msg.query)
		}
		return m, nil

	case prefetchTickMsg:
		if msg.gen == m.selectionGen {
			m.cache.Prefetch(msg.videoID)
		}
		return m, nil

	case trackReadyMsg:
		if msg.err != nil {
			errMsg := msg.err.Error()
			if strings.Contains(errMsg, "Sign in to confirm") || strings.Contains(errMsg, "bot") {
				m.statusMsg = "YouTube bot detection. Please run: yt-dlp --cookies-from-browser firefox --cookies cookies.txt"
			} else if strings.Contains(errMsg, "429") {
				m.statusMsg = "Too many requests. Please wait a few minutes and try again."
			} else {
				m.statusMsg = fmt.Sprintf("Error: %v", msg.err)
			}
			m.isPlaying = false
		} else {
			m.isPlaying = true
			m.currentTrack = &msg.track
			m.startTime = time.Now()
			m.elapsedTime = 0
			m.totalMinutes += int(msg.track.Duration / 60)
			m.statusMsg = fmt.Sprintf("Playing: %s", msg.track.Title)
		}
		return m, nil
	}

	if m.isSearching {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "enter":
				query := strings.TrimSpace(m.searchInput.Value())
				m.isSearching = false
				m.searchInput.Blur()
				if query == "" {
					m.statusMsg = "Search cancelled."
					return m, nil
				}
				if len(m.searchResults) > 0 {
					m.statusMsg = fmt.Sprintf("%d results found", len(m.searchResults))
					return m, nil
				}
				m.statusMsg = "Searching with yt-dlp..."
				return m, performSearch(query)

			case "esc":
				m.isSearching = false
				m.searchInput.Blur()
				m.statusMsg = "Search cancelled."
				return m, nil
			}
		}

		previousValue := m.searchInput.Value()
		var updateCmd tea.Cmd
		m.searchInput, updateCmd = m.searchInput.Update(msg)

		if m.searchInput.Value() != previousValue {
			query := strings.TrimSpace(m.searchInput.Value())
			if len(query) >= 2 {
				m.searchGen++
				m.statusMsg = "Searching with yt-dlp..."
				return m, tea.Batch(updateCmd, debounceSearch(m.searchGen, query))
			}
		}

		return m, updateCmd
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "q", "ctrl+c":
			if m.engine != nil {
				m.engine.Stop()
			}
			return m, tea.Quit

		case "/":
			m.isSearching = true
			m.searchInput.Focus()
			m.searchInput.SetValue("")
			m.statusMsg = "Search: type your query"
			return m, nil

		case "j", "down":
			if m.cursor < len(m.searchResults)-1 {
				m.cursor++
			}
			if len(m.searchResults) > 0 {
				m.selectionGen++
				return m, debouncePrefetch(m.selectionGen, m.searchResults[m.cursor].ID)
			}
			return m, nil

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
			if len(m.searchResults) > 0 {
				m.selectionGen++
				return m, debouncePrefetch(m.selectionGen, m.searchResults[m.cursor].ID)
			}
			return m, nil

		case " ":
			if m.engine != nil && m.currentTrack != nil {
				err := m.engine.TogglePause()
				if err != nil {
					m.statusMsg = fmt.Sprintf("Error: %v", err)
				} else {
					m.isPlaying = !m.isPlaying
					if m.isPlaying {
						m.startTime = time.Now().Add(-m.elapsedTime)
						m.statusMsg = "Playing"
					} else {
						m.statusMsg = "Paused"
					}
				}
			} else {
				m.statusMsg = "No track playing"
			}
			return m, nil

		case "enter":
			if len(m.searchResults) > 0 && m.cursor < len(m.searchResults) {
				selected := m.searchResults[m.cursor]
				m.currentTrack = &selected
				m.statusMsg = fmt.Sprintf("Loading: %s...", selected.Title)
				return m, fetchTrack(m.engine, m.cache, selected)
			}
			return m, nil
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

	header := HeaderStyle.Width(m.width - 4).Render(
		fmt.Sprintf("Nelo TunePipe %s %s SCROBBLES: %s %s",
			Divider(),
			CounterStyle.Render(formatDuration(m.totalMinutes)),
			Divider(),
			CounterStyle.Render(formatScrobbles(m.scrobbles)),
		),
	)

	var searchContent string
	if m.isSearching {
		searchBox := SearchInputStyle.Width(mainWidth - 4).Render(m.searchInput.View())
		searchContent = SearchModeStyle.Render("SEARCH") + "\n\n" + searchBox
		searchContent += "\n\n" + MutedStyle.Render("Enter to confirm • Esc to cancel")
		if len(m.searchResults) > 0 {
			searchContent += "\n\n" + renderResultsList(m)
		}
	} else if len(m.searchResults) == 0 {
		searchContent = MutedStyle.Render("Press '/' to search for music")
	} else {
		searchContent = renderResultsList(m)
	}

	mainPane := MainPaneStyle.
		Width(mainWidth).
		Height(10).
		Render("MUSIC\n\n" + searchContent)

	sidePane := SidebarStyle.
		Width(30).
		Height(10).
		Render(
			"CONTROLS\n\n" +
				"[/] Search\n" +
				"[j/k] Navigate\n" +
				"[Enter] Play\n" +
				"[Space] Pause\n" +
				"[q] Quit\n\n" +
				TagStyle.Render("#tunepipe") + " " +
				TagStyle.Render("#music"),
		)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, mainPane, sidePane)

	nowPlaying := "No track playing"
	elapsedStr := ""
	if m.currentTrack != nil {
		if m.isPlaying {
			elapsedStr = fmt.Sprintf("%02d:%02d / %02d:%02d",
				int(m.elapsedTime.Minutes()),
				int(m.elapsedTime.Seconds())%60,
				int(m.currentTrack.Duration/60),
				int(m.currentTrack.Duration)%60,
			)
		} else {
			elapsedStr = fmt.Sprintf("%02d:%02d / %02d:%02d  [PAUSED]",
				int(m.elapsedTime.Minutes()),
				int(m.elapsedTime.Seconds())%60,
				int(m.currentTrack.Duration/60),
				int(m.currentTrack.Duration)%60,
			)
		}
		nowPlaying = fmt.Sprintf("%s %s %s",
			TitleStyle.Render(m.currentTrack.Title),
			MutedStyle.Render("•"),
			MutedStyle.Render(m.currentTrack.Uploader),
		)
	}

	statusStyle := StatusPausedStyle
	if m.isPlaying {
		statusStyle = StatusPlayingStyle
	} else if strings.Contains(m.statusMsg, "bot") || strings.Contains(m.statusMsg, "429") {
		statusStyle = StatusErrorStyle
	}

	progress := ""
	if m.currentTrack != nil && m.currentTrack.Duration > 0 {
		pct := int((m.elapsedTime.Seconds() / m.currentTrack.Duration) * 20)
		if pct > 20 {
			pct = 20
		}
		bar := strings.Repeat("█", pct) + strings.Repeat("░", 20-pct)
		progress = ProgressBarStyle.Render(bar)
	}

	playerBar := PlayerBarStyle.
		Width(playerWidth).
		Render(
			fmt.Sprintf("%s\n%s %s\n%s",
				statusStyle.Render(nowPlaying),
				CounterStyle.Render(elapsedStr),
				Divider(),
				progress,
			) +
				"\n" + MutedStyle.Render(m.statusMsg),
		)

	return lipgloss.JoinVertical(lipgloss.Left, header, topRow, playerBar)
}

func renderResultsList(m Model) string {
	var content string
	for i, track := range m.searchResults {
		cursorStr := " "
		item := fmt.Sprintf("%2d. %s %s %s",
			i+1,
			TitleStyle.Render(truncate(track.Title, 30)),
			MutedStyle.Render("•"),
			MutedStyle.Render(truncate(track.Uploader, 20)),
		)
		if m.cursor == i {
			cursorStr = ">"
			itemStr := fmt.Sprintf("%s %s", cursorStr, item)
			content += SelectedItemStyle.Render(itemStr) + "\n"
		} else {
			itemStr := fmt.Sprintf("%s %s", cursorStr, item)
			content += NormalItemStyle.Render(itemStr) + "\n"
		}
	}
	return content
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen < 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func formatDuration(minutes int) string {
	hours := minutes / 60
	mins := minutes % 60
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	}
	return fmt.Sprintf("%dm", mins)
}

func formatScrobbles(scrobbles int) string {
	if scrobbles >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(scrobbles)/1000000)
	}
	if scrobbles >= 1000 {
		return fmt.Sprintf("%.1fK", float64(scrobbles)/1000)
	}
	return strconv.Itoa(scrobbles)
}
