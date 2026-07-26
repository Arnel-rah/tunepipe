package ui

import (
	"fmt"
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
	statusKind    string
	width         int
	height        int

	totalMinutes int
	elapsedTime  time.Duration
	scrobbles    int
	tickCount    int

	confirmQuit bool
}

func NewModel(engine *player.Engine) Model {
	ti := textinput.New()
	ti.Placeholder = "search for a song..."
	ti.CharLimit = 156
	ti.Prompt = ""
	ti.TextStyle = lipgloss.NewStyle().Foreground(TextColor)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(TextMutedColor)

	return Model{
		engine:       engine,
		cache:        ytdlp.NewURLCache(5 * time.Minute),
		searchInput:  ti,
		statusMsg:    "ready",
		statusKind:   "idle",
		width:        80,
		height:       24,
		isPlaying:    false,
		totalMinutes: 0,
		scrobbles:    0,
		confirmQuit:  false,
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
		m.tickCount++
		if m.isPlaying && m.currentTrack != nil {
			m.elapsedTime += time.Second

			if int(m.elapsedTime.Seconds())%10 == 0 {
				m.scrobbles++
			}

			if m.elapsedTime.Seconds() >= m.currentTrack.Duration {
				m.isPlaying = false
				m.statusMsg = "finished"
				m.statusKind = "idle"
				m.elapsedTime = time.Duration(m.currentTrack.Duration) * time.Second
			}
		}
		return m, tickCmd()

	case searchResultsMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("error: %v", msg.err)
			m.statusKind = "error"
			return m, nil
		}
		m.searchResults = msg.tracks
		m.cursor = 0
		if len(msg.tracks) > 0 {
			m.statusMsg = fmt.Sprintf("%d results", len(msg.tracks))
			m.statusKind = "idle"
		} else {
			m.statusMsg = "no results"
			m.statusKind = "idle"
		}

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
				m.statusMsg = "bot detection - export cookies"
			} else if strings.Contains(errMsg, "429") {
				m.statusMsg = "too many requests, wait..."
			} else {
				m.statusMsg = fmt.Sprintf("error: %v", msg.err)
			}
			m.statusKind = "error"
			m.isPlaying = false
		} else {
			m.isPlaying = true
			m.currentTrack = &msg.track
			m.elapsedTime = 0
			m.totalMinutes += int(msg.track.Duration / 60)
			m.statusMsg = "playing"
			m.statusKind = "playing"
		}
		return m, nil
	}

	if m.confirmQuit {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "y", "Y":
				if m.engine != nil {
					m.engine.Stop()
				}
				return m, tea.Quit
			case "n", "N", "esc":
				m.confirmQuit = false
				m.statusMsg = "ready"
				m.statusKind = "idle"
				return m, nil
			}
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
					m.statusMsg = "search cancelled"
					m.statusKind = "idle"
					return m, nil
				}
				m.statusMsg = "searching..."
				m.statusKind = "idle"
				return m, performSearch(query)

			case "esc":
				m.isSearching = false
				m.searchInput.Blur()
				m.statusMsg = "search cancelled"
				m.statusKind = "idle"
				return m, nil
			}
		}

		var updateCmd tea.Cmd
		m.searchInput, updateCmd = m.searchInput.Update(msg)
		return m, updateCmd
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "q", "ctrl+c":
			m.confirmQuit = true
			return m, nil

		case "/":
			m.isSearching = true
			m.searchInput.Focus()
			m.searchInput.SetValue("")
			m.statusMsg = "search..."
			m.statusKind = "idle"
			return m, nil

		case "j", "down":
			if m.cursor < len(m.searchResults)-1 {
				m.cursor++
				if len(m.searchResults) > 0 {
					m.selectionGen++
					return m, debouncePrefetch(m.selectionGen, m.searchResults[m.cursor].ID)
				}
			}

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
				if len(m.searchResults) > 0 {
					m.selectionGen++
					return m, debouncePrefetch(m.selectionGen, m.searchResults[m.cursor].ID)
				}
			}

		case " ":
			if m.engine != nil && m.currentTrack != nil {
				if m.elapsedTime.Seconds() >= m.currentTrack.Duration {
					m.isPlaying = false
					m.statusMsg = "finished"
					m.statusKind = "idle"
					return m, nil
				}
				err := m.engine.TogglePause()
				if err != nil {
					m.statusMsg = fmt.Sprintf("error: %v", err)
					m.statusKind = "error"
				} else {
					m.isPlaying = !m.isPlaying
					if m.isPlaying {
						m.statusMsg = "playing"
						m.statusKind = "playing"
					} else {
						m.statusMsg = "paused"
						m.statusKind = "paused"
					}
				}
			} else {
				m.statusMsg = "no track loaded"
				m.statusKind = "idle"
			}
			return m, nil

		case "enter":
			if len(m.searchResults) > 0 && m.cursor < len(m.searchResults) {
				selected := m.searchResults[m.cursor]
				m.currentTrack = &selected
				m.elapsedTime = 0
				m.statusMsg = fmt.Sprintf("loading: %s", selected.Title)
				m.statusKind = "idle"
				return m, fetchTrack(m.engine, m.cache, selected)
			}
		}
	}

	return m, nil
}
