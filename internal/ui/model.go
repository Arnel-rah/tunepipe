package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Arnel-rah/tunepipe/internal/notify"
	"github.com/Arnel-rah/tunepipe/internal/player"
	"github.com/Arnel-rah/tunepipe/internal/ytdlp"
	"github.com/Arnel-rah/tunepipe/internal/notify"

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

const (
	prefetchDelay = 150 * time.Millisecond
	searchDelay   = 500 * time.Millisecond
	cacheTTL      = 5 * time.Minute
)

var spinnerFrames = [...]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

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

	queue       []ytdlp.Track
	queueSet    map[string]struct{}
	showQueue   bool
	queueCursor int

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
		engine:      engine,
		cache:       ytdlp.NewURLCache(cacheTTL),
		searchInput: ti,
		statusMsg:   "ready",
		statusKind:  "idle",
		width:       80,
		height:      24,
		queueSet:    make(map[string]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tickCmd())
}

func (m Model) SpinnerFrame() string {
	return spinnerFrames[m.tickCount%len(spinnerFrames)]
}

func (m Model) IsQueued(trackID string) bool {
	_, ok := m.queueSet[trackID]
	return ok
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func debouncePrefetch(gen int, videoID string) tea.Cmd {
	return tea.Tick(prefetchDelay, func(time.Time) tea.Msg {
		return prefetchTickMsg{gen: gen, videoID: videoID}
	})
}

func debounceSearch(gen int, query string) tea.Cmd {
	return tea.Tick(searchDelay, func(time.Time) tea.Msg {
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
		// Prefer local cached file for offline playback
		if local := ytdlp.CachedFilePath(track.ID); local != "" {
			if err := engine.PlayURL(local); err != nil {
				return trackReadyMsg{track: track, err: err}
			}
			return trackReadyMsg{track: track, directURL: local}
		}

		directURL, err := cache.Get(track.ID)
		if err != nil {
			return trackReadyMsg{track: track, err: err}
		}
		if err := engine.PlayURL(directURL); err != nil {
			return trackReadyMsg{track: track, err: err}
		}
		return trackReadyMsg{track: track, directURL: directURL}
	}
}

func playTrack(m Model, track ytdlp.Track) (Model, tea.Cmd) {
	m.currentTrack = &track
	m.elapsedTime = 0
	m.isPlaying = false
	m.isLoading = true
	m.statusMsg = "loading..."
	m.statusKind = "idle"

	if m.isSearching {
		m.isSearching = false
		m.searchInput.Blur()
	}

	return m, fetchTrack(m.engine, m.cache, track)
}

func playTrackAt(m Model, i int) (Model, tea.Cmd) {
	if i < 0 || i >= len(m.searchResults) {
		return m, nil
	}
	m.cursor = i
	return playTrack(m, m.searchResults[i])
}

func dequeueNext(m Model) (ytdlp.Track, Model, bool) {
	if len(m.queue) == 0 {
		return ytdlp.Track{}, m, false
	}
	next := m.queue[0]
	m.queue = m.queue[1:]
	delete(m.queueSet, next.ID)
	return next, m, true
}

func toggleQueue(m Model, track ytdlp.Track) (Model, bool) {
	if _, exists := m.queueSet[track.ID]; exists {
		// remove
		for i, t := range m.queue {
			if t.ID == track.ID {
				m.queue = append(m.queue[:i], m.queue[i+1:]...)
				break
			}
		}
		delete(m.queueSet, track.ID)
		return m, false
	}

	m.queue = append(m.queue, track)
	m.queueSet[track.ID] = struct{}{}
	return m, true
}

func playNextAuto(m Model) (Model, tea.Cmd) {
	if next, updated, ok := dequeueNext(m); ok {
		return playTrack(updated, next)
	}
	if m.cursor < len(m.searchResults)-1 {
		return playTrackAt(m, m.cursor+1)
	}
	return m, nil
}

func setStatus(m *Model, msg, kind string) {
	m.statusMsg = msg
	m.statusKind = kind
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
			secs := int(m.elapsedTime.Seconds())

			if secs%10 == 0 {
				m.scrobbles++
			}

			if float64(secs) >= m.currentTrack.Duration {
				m.elapsedTime = time.Duration(m.currentTrack.Duration) * time.Second
				m.isPlaying = false
				setStatus(&m, "finished", "idle")

				nextModel, cmd := playNextAuto(m)
				if cmd != nil {
					return nextModel, tea.Batch(tickCmd(), cmd)
				}
			}
		}
		return m, tickCmd()

	case searchResultsMsg:
		m.searchPending = false
		if msg.err != nil {
			setStatus(&m, fmt.Sprintf("error: %v", msg.err), "error")
			return m, nil
		}

		m.searchResults = msg.tracks
		m.cursor = 0

		if n := len(msg.tracks); n > 0 {
			setStatus(&m, fmt.Sprintf("%d results", n), "idle")
			m.cache.PrefetchWindow(msg.tracks, m.cursor)
		} else {
			setStatus(&m, "no results", "idle")
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
			switch {
			case strings.Contains(errMsg, "Sign in to confirm") || strings.Contains(errMsg, "bot"):
				setStatus(&m, "bot detection - export cookies", "error")
			case strings.Contains(errMsg, "429"):
				setStatus(&m, "too many requests, wait...", "error")
			default:
				setStatus(&m, fmt.Sprintf("error: %v", msg.err), "error")
			}
			m.isPlaying = false
			m.isLoading = false
			m.elapsedTime = 0
			return m, nil
		}

		m.isPlaying = true
		m.isLoading = false
		m.currentTrack = &msg.track
		m.elapsedTime = 0
		m.totalMinutes += int(msg.track.Duration / 60)
		setStatus(&m, "playing", "playing")

		m.cache.PrefetchWindow(m.searchResults, m.cursor)
		if len(m.queue) > 0 {
			m.cache.PrefetchWindow(m.queue, 0)
		}

		if count := ytdlp.IncrementPlayCount(msg.track.ID); count >= 3 {
			id := msg.track.ID
			go func() { _, _ = ytdlp.EnsureDownloaded(id) }()
		}

		t := msg.track
		go func() { _ = notify.NotifyTrack(t) }()

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
				setStatus(&m, "ready", "idle")
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
					setStatus(&m, "search cancelled", "idle")
					m.isSearching = false
					m.searchInput.Blur()
					return m, nil
				}
				m.isSearching = false
				m.searchPending = true
				m.searchInput.Blur()
				m.searchGen++
				setStatus(&m, "searching...", "idle")
				return m, performSearch(query)

			case "esc":
				m.isSearching = false
				setStatus(&m, "search cancelled", "idle")
				m.searchInput.Blur()
				return m, nil
			}
		}

		prev := m.searchInput.Value()
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)

		if newVal := m.searchInput.Value(); newVal != prev {
			query := strings.TrimSpace(newVal)
			if query == "" {
				setStatus(&m, "search...", "idle")
				return m, cmd
			}
			m.searchGen++
			return m, tea.Batch(cmd, debounceSearch(m.searchGen, query))
		}
		return m, cmd
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch keyMsg.String() {
	case "q", "ctrl+c":
		m.confirmQuit = true
		return m, nil

	case "/":
		m.showQueue = false
		m.isSearching = true
		m.searchInput.Focus()
		m.searchInput.SetValue("")
		setStatus(&m, "search...", "idle")
		return m, nil

	case "Q":
		if m.showQueue || len(m.searchResults) > 0 || len(m.queue) > 0 {
			m.showQueue = !m.showQueue
			if m.showQueue && m.queueCursor >= len(m.queue) {
				m.queueCursor = 0
			}
		}
		return m, nil

	case "j", "down":
		if m.showQueue {
			if m.queueCursor < len(m.queue)-1 {
				m.queueCursor++
			}
			return m, nil
		}
		if m.cursor < len(m.searchResults)-1 {
			m.cursor++
			m.selectionGen++
			return m, debouncePrefetch(m.selectionGen, m.searchResults[m.cursor].ID)
		}
		return m, nil

	case "k", "up":
		if m.showQueue {
			if m.queueCursor > 0 {
				m.queueCursor--
			}
			return m, nil
		}
		if m.cursor > 0 {
			m.cursor--
			m.selectionGen++
			return m, debouncePrefetch(m.selectionGen, m.searchResults[m.cursor].ID)
		}
		return m, nil

	case "a":
		if !m.showQueue && m.cursor < len(m.searchResults) {
			track := m.searchResults[m.cursor]
			var added bool
			m, added = toggleQueue(m, track)
			if added {
				setStatus(&m, fmt.Sprintf("added to queue (%d)", len(m.queue)), "idle")
			} else {
				setStatus(&m, fmt.Sprintf("removed from queue (%d)", len(m.queue)), "idle")
			}
		}
		return m, nil

	case "A":
		if len(m.queue) > 0 {
			m.queue = nil
			m.queueSet = make(map[string]struct{})
			m.queueCursor = 0
			m.showQueue = false
			setStatus(&m, "queue cleared", "idle")
		}
		return m, nil

	case "x":
		if m.showQueue && m.queueCursor < len(m.queue) {
			id := m.queue[m.queueCursor].ID
			m.queue = append(m.queue[:m.queueCursor], m.queue[m.queueCursor+1:]...)
			delete(m.queueSet, id)
			if m.queueCursor >= len(m.queue) && m.queueCursor > 0 {
				m.queueCursor--
			}
			setStatus(&m, fmt.Sprintf("removed from queue (%d)", len(m.queue)), "idle")
		}
		return m, nil

	case " ":
		if m.engine == nil || m.currentTrack == nil {
			setStatus(&m, "no track loaded", "idle")
			return m, nil
		}
		if m.isLoading {
			setStatus(&m, "loading...", "idle")
			return m, nil
		}
		if m.elapsedTime.Seconds() >= m.currentTrack.Duration {
			m.isPlaying = false
			setStatus(&m, "finished", "idle")
			return m, nil
		}

		if m.isPlaying {
			if err := m.engine.Pause(); err != nil {
				setStatus(&m, fmt.Sprintf("error: %v", err), "error")
				return m, nil
			}
			m.isPlaying = false
			setStatus(&m, "paused", "paused")
		} else {
			if err := m.engine.Resume(); err != nil {
				setStatus(&m, fmt.Sprintf("error: %v", err), "error")
				return m, nil
			}
			m.isPlaying = true
			setStatus(&m, "playing", "playing")
		}
		return m, nil

	case "n":
		if len(m.queue) > 0 {
			return playNextAuto(m)
		}
		setStatus(&m, "queue empty", "idle")
		return m, nil

	case "enter":
		if m.showQueue {
			if m.queueCursor < len(m.queue) {
				track := m.queue[m.queueCursor]
				id := track.ID
				m.queue = append(m.queue[:m.queueCursor], m.queue[m.queueCursor+1:]...)
				delete(m.queueSet, id)
				if m.queueCursor >= len(m.queue) && m.queueCursor > 0 {
					m.queueCursor--
				}
				return playTrack(m, track)
			}
			return m, nil
		}

		if m.cursor < len(m.searchResults) {
			if m.currentTrack != nil && m.searchResults[m.cursor].ID == m.currentTrack.ID {
				return m, nil
			}
			return playTrackAt(m, m.cursor)
		}
		return m, nil
	}

	return m, nil
}
