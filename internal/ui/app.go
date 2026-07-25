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
	ti.TextStyle = lipgloss.NewStyle().Foreground(TextPrimary)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(TextDim)

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
				m.statusMsg = "quit cancelled"
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
				if len(m.searchResults) > 0 {
					m.statusMsg = fmt.Sprintf("%d results", len(m.searchResults))
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

		previousValue := m.searchInput.Value()
		var updateCmd tea.Cmd
		m.searchInput, updateCmd = m.searchInput.Update(msg)

		if m.searchInput.Value() != previousValue {
			query := strings.TrimSpace(m.searchInput.Value())
			if len(query) >= 2 {
				m.searchGen++
				m.statusMsg = "searching..."
				m.statusKind = "idle"
				return m, tea.Batch(updateCmd, debounceSearch(m.searchGen, query))
			}
		}

		return m, updateCmd
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "q", "ctrl+c":
			m.confirmQuit = true
			m.statusMsg = "quit? (y/n)"
			m.statusKind = "idle"
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

func (m Model) View() string {
	if m.width == 0 {
		m.width = 80
	}

	width := m.width - 2
	if width < 40 {
		width = 40
	}

	header := renderHeader(m, width)
	nowPlaying := renderNowPlaying(m, width)
	search := renderSearch(m, width)
	results := renderResults(m, width)
	progress := renderProgress(m, width)
	controls := renderControls(m, width)
	footer := renderFooter(m, width)

	body := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		nowPlaying,
		"",
		search,
		"",
		results,
		"",
		progress,
		"",
		controls,
		"",
		footer,
	)

	if m.confirmQuit {
		confirmBox := ConfirmStyle.Render(
			fmt.Sprintf("%s %s %s",
				AccentTextStyle.Render(IconBullet),
				TitleStyle.Render("Are you sure you want to quit?"),
				MutedStyle.Render("(y/n)"),
			),
		)
		body = lipgloss.JoinVertical(lipgloss.Center, body, "", confirmBox)
	}

	return body
}

func renderHeader(m Model, width int) string {
	logo := LogoStyle.Render(IconDiamond + " tunepipe")

	stats := fmt.Sprintf("%s %dh  %s %d",
		StatLabelStyle.Render("listened"),
		m.totalMinutes/60,
		StatLabelStyle.Render("scrobbles"),
		m.scrobbles,
	)

	gap := width - lipgloss.Width(logo) - lipgloss.Width(stats) - 2
	if gap < 1 {
		gap = 1
	}

	return fmt.Sprintf("%s%s%s", logo, strings.Repeat(" ", gap), stats)
}

func renderNowPlaying(m Model, width int) string {
	if m.currentTrack == nil {
		return MutedStyleDim.Render(IconDiamond + " nothing playing " + IconDiamond)
	}

	var status string
	var style lipgloss.Style
	if m.isPlaying {
		status = IconPlay
		style = StatusTextPlaying
	} else if m.elapsedTime.Seconds() >= m.currentTrack.Duration {
		status = IconStop
		style = StatusTextIdle
	} else {
		status = IconPause
		style = StatusTextPaused
	}

	title := TitleStyle.Render(truncate(m.currentTrack.Title, width-20))
	artist := MutedStyle.Render(truncate(m.currentTrack.Uploader, width-20))

	return fmt.Sprintf("%s %s %s %s",
		style.Render(status),
		DividerStyle.Render(IconArrow),
		title,
		MutedStyleDim.Render("by "+artist),
	)
}

func renderSearch(m Model, width int) string {
	if m.isSearching {
		return SearchInputFocusedStyle.Width(width).Render(m.searchInput.View())
	}

	if len(m.searchResults) > 0 {
		return MutedStyleDim.Render(fmt.Sprintf("%s %d results • press / to search",
			IconDiamond, len(m.searchResults)))
	}

	return MutedStyleDim.Render(IconDiamond + " press / to search")
}

func renderResults(m Model, width int) string {
	if len(m.searchResults) == 0 {
		return ""
	}

	var lines []string
	maxItems := 8

	for i := 0; i < len(m.searchResults) && i < maxItems; i++ {
		track := m.searchResults[i]
		line := renderResultLine(track, i, i == m.cursor, width)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func renderResultLine(track ytdlp.Track, index int, selected bool, width int) string {
	num := fmt.Sprintf("%2d.", index+1)
	title := truncate(track.Title, width-30)
	artist := truncate(track.Uploader, 20)

	if selected {
		return RowSelectedStyle.Render(fmt.Sprintf("%s %s %s  %s",
			IconArrow, num, title, artist))
	}
	return RowStyle.Render(fmt.Sprintf("  %s %s  %s", num, title, artist))
}

func renderProgress(m Model, width int) string {
	if m.currentTrack == nil || m.currentTrack.Duration <= 0 {
		return ""
	}

	barWidth := width - 16
	if barWidth < 10 {
		barWidth = 10
	}

	var pct float64
	if m.elapsedTime.Seconds() >= m.currentTrack.Duration {
		pct = 1.0
	} else if m.currentTrack.Duration > 0 {
		pct = m.elapsedTime.Seconds() / m.currentTrack.Duration
		if pct > 1 {
			pct = 1
		}
		if pct < 0 {
			pct = 0
		}
	}

	filled := int(pct * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	if filled < 0 {
		filled = 0
	}

	elapsed := fmt.Sprintf("%02d:%02d", int(m.elapsedTime.Minutes()), int(m.elapsedTime.Seconds())%60)
	total := fmt.Sprintf("%02d:%02d", int(m.currentTrack.Duration/60), int(m.currentTrack.Duration)%60)

	isFinished := m.elapsedTime.Seconds() >= m.currentTrack.Duration

	var bar string
	if isFinished {
		bar = MutedStyleDim.Render(strings.Repeat(CharProgress, barWidth))
	} else if filled <= 0 {
		bar = MutedStyleDim.Render(strings.Repeat(CharTrack, barWidth))
	} else if filled >= barWidth {
		bar = ProgressBarStyle.Render(strings.Repeat(CharProgress, barWidth))
	} else {
		bar = ProgressBarStyle.Render(strings.Repeat(CharProgress, filled)) +
			MutedStyleDim.Render(strings.Repeat(CharTrack, barWidth-filled))
	}

	if isFinished {
		return fmt.Sprintf("%s %s %s", MutedStyleDim.Render(elapsed), MutedStyleDim.Render(bar), MutedStyleDim.Render(total))
	}

	return fmt.Sprintf("%s %s %s", MutedStyle.Render(elapsed), bar, MutedStyle.Render(total))
}

func renderControls(m Model, width int) string {
	keys := []struct {
		key   string
		label string
	}{
		{"space", "play/pause"},
		{"/", "search"},
		{"j/k", "navigate"},
		{"enter", "select"},
		{"q", "quit"},
	}

	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s %s",
			KeyHintKeyStyle.Render(k.key),
			KeyHintTextStyle.Render(k.label),
		))
	}

	line := strings.Join(parts, "  "+IconBullet+"  ")
	return MutedStyleDim.Render(line)
}

func renderFooter(m Model, width int) string {
	status := m.statusMsg
	if status == "" {
		status = "ready"
	}

	var style lipgloss.Style
	var icon string
	switch m.statusKind {
	case "playing":
		style = StatusTextPlaying
		icon = IconPlay
	case "paused":
		style = StatusTextPaused
		icon = IconPause
	case "error":
		style = StatusTextError
		icon = IconCross
	default:
		style = StatusTextIdle
		icon = IconBullet
	}

	statusText := style.Render(icon + " " + status)

	keys := MutedStyleDim.Render("space · / · j/k · enter · q")

	gap := width - lipgloss.Width(statusText) - lipgloss.Width(keys) - 2
	if gap < 1 {
		gap = 1
	}

	return fmt.Sprintf("%s%s%s", statusText, strings.Repeat(" ", gap), keys)
}

func truncate(s string, maxLen int) string {
	if maxLen < 1 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen < 3 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-3]) + "…"
}
