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
	searchPending bool
	searchResults []ytdlp.Track
	cursor        int
	selectionGen  int
	searchGen     int
	currentTrack  *ytdlp.Track
	isPlaying     bool
	isLoading     bool
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
		isLoading:    false,
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

func (m Model) SpinnerFrame() string {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	return frames[m.tickCount%len(frames)]
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
		results, err := ytdlp.SearchWithCache(query, 10)
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

func playTrackAt(m Model, i int) (Model, tea.Cmd) {
	if i < 0 || i >= len(m.searchResults) {
		return m, nil
	}
	selected := m.searchResults[i]
	m.cursor = i
	m.currentTrack = &selected
	m.elapsedTime = 0
	m.isPlaying = false
	m.isLoading = true
	m.statusMsg = "loading..."
	m.statusKind = "idle"

	if m.isSearching {
		m.isSearching = false
		m.searchInput.Blur()
	}

	return m, fetchTrack(m.engine, m.cache, selected)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		m.tickCount++
		if m.isPlaying && m.currentTrack != nil && !m.isLoading {
			m.elapsedTime += time.Second

			if int(m.elapsedTime.Seconds())%10 == 0 {
				m.scrobbles++
			}

			if m.elapsedTime.Seconds() >= m.currentTrack.Duration {
				m.elapsedTime = time.Duration(m.currentTrack.Duration) * time.Second
				m.isPlaying = false
				m.statusMsg = "finished"
				m.statusKind = "idle"

				if m.cursor < len(m.searchResults)-1 {
					nextTrack, cmd := playTrackAt(m, m.cursor+1)
					return nextTrack, tea.Batch(tickCmd(), cmd)
				}
			}
		}
		return m, tickCmd()

	case searchResultsMsg:
		m.searchPending = false
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

		m.cache.PrefetchWindow(msg.tracks, m.cursor)
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
			m.isLoading = false
			m.elapsedTime = 0
		} else {
			m.isPlaying = true
			m.isLoading = false
			m.currentTrack = &msg.track
			m.elapsedTime = 0
			m.totalMinutes += int(msg.track.Duration / 60)
			m.statusMsg = "playing"
			m.statusKind = "playing"

			m.cache.PrefetchWindow(m.searchResults, m.cursor)
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
				if query == "" {
					m.statusMsg = "search cancelled"
					m.statusKind = "idle"
					m.isSearching = false
					m.searchInput.Blur()
					return m, nil
				}
				m.isSearching = false
				m.searchPending = true
				m.searchInput.Blur()
				m.searchGen++
				m.statusMsg = "searching..."
				m.statusKind = "idle"
				return m, performSearch(query)

			case "esc":
				m.isSearching = false
				m.statusMsg = "search cancelled"
				m.statusKind = "idle"
				m.searchInput.Blur()
				return m, nil
			}
		}

		prevValue := m.searchInput.Value()
		var updateCmd tea.Cmd
		m.searchInput, updateCmd = m.searchInput.Update(msg)

		if newValue := m.searchInput.Value(); newValue != prevValue {
			query := strings.TrimSpace(newValue)
			if query == "" {
				m.statusMsg = "search..."
				m.statusKind = "idle"
				return m, updateCmd
			}
			m.searchGen++
			return m, tea.Batch(updateCmd, debounceSearch(m.searchGen, query))
		}

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
			return m, nil

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
				if len(m.searchResults) > 0 {
					m.selectionGen++
					return m, debouncePrefetch(m.selectionGen, m.searchResults[m.cursor].ID)
				}
			}
			return m, nil

		case " ":
			if m.engine == nil || m.currentTrack == nil {
				m.statusMsg = "no track loaded"
				m.statusKind = "idle"
				return m, nil
			}

			if m.isLoading {
				m.statusMsg = "loading..."
				m.statusKind = "idle"
				return m, nil
			}

			if m.elapsedTime.Seconds() >= m.currentTrack.Duration {
				m.isPlaying = false
				m.statusMsg = "finished"
				m.statusKind = "idle"
				return m, nil
			}

			if m.isPlaying {
				if err := m.engine.Pause(); err != nil {
					m.statusMsg = fmt.Sprintf("error: %v", err)
					m.statusKind = "error"
					return m, nil
				}
				m.isPlaying = false
				m.statusMsg = "paused"
				m.statusKind = "paused"
			} else {
				if err := m.engine.Resume(); err != nil {
					m.statusMsg = fmt.Sprintf("error: %v", err)
					m.statusKind = "error"
					return m, nil
				}
				m.isPlaying = true
				m.statusMsg = "playing"
				m.statusKind = "playing"
			}
			return m, nil

		case "enter":
			if len(m.searchResults) > 0 && m.cursor < len(m.searchResults) {
				if m.currentTrack != nil && m.searchResults[m.cursor].ID == m.currentTrack.ID {
					return m, nil
				}
				return playTrackAt(m, m.cursor)
			}
			return m, nil
		}
	}

	return m, nil
}
