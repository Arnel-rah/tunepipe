package lyrics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Line struct {
	At   time.Duration
	Text string
}

type Client struct {
	httpClient *http.Client
}

type response struct {
	SyncedLyrics string `json:"syncedLyrics"`
	PlainLyrics  string `json:"plainLyrics"`
}

type searchResult struct {
	SyncedLyrics string `json:"syncedLyrics"`
	PlainLyrics  string `json:"plainLyrics"`
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 8 * time.Second}
	}
	return &Client{httpClient: httpClient}
}

func (c *Client) Get(ctx context.Context, artist, title string) ([]Line, error) {
	artist = cleanMetadata(artist)
	title = cleanMetadata(title)
	lines, err := c.get(ctx, artist, title)
	if err == nil {
		return lines, nil
	}

	query := url.Values{"q": {strings.TrimSpace(artist + " " + title)}}
	req, requestErr := http.NewRequestWithContext(ctx, http.MethodGet, "https://lrclib.net/api/search?"+query.Encode(), nil)
	if requestErr != nil {
		return nil, requestErr
	}
	req.Header.Set("User-Agent", "tunepipe/lyrics")
	resp, requestErr := c.httpClient.Do(req)
	if requestErr != nil {
		return nil, fmt.Errorf("lyrics search failed: %w", requestErr)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, err
	}
	var results []searchResult
	if requestErr := json.NewDecoder(resp.Body).Decode(&results); requestErr != nil {
		return nil, fmt.Errorf("decode lyrics search response: %w", requestErr)
	}
	for _, result := range results {
		if result.SyncedLyrics != "" {
			return parseSynced(result.SyncedLyrics)
		}
		if result.PlainLyrics != "" {
			return plainLines(result.PlainLyrics), nil
		}
	}
	return nil, err
}

func (c *Client) get(ctx context.Context, artist, title string) ([]Line, error) {
	query := url.Values{"artist_name": {artist}, "track_name": {title}}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://lrclib.net/api/get?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "tunepipe/lyrics")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("lyrics request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("lyrics not found")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("lyrics service returned HTTP %d", resp.StatusCode)
	}

	var payload response
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode lyrics response: %w", err)
	}
	if payload.SyncedLyrics != "" {
		return parseSynced(payload.SyncedLyrics)
	}
	if payload.PlainLyrics != "" {
		return plainLines(payload.PlainLyrics), nil
	}
	return nil, fmt.Errorf("lyrics response is empty")
}

func cleanMetadata(value string) string {
	for _, separator := range []string{"[", "(", "|", " - "} {
		if index := strings.Index(value, separator); index > 0 {
			value = value[:index]
		}
	}
	return strings.TrimSpace(value)
}

func parseSynced(raw string) ([]Line, error) {
	var lines []Line
	for _, rawLine := range strings.Split(raw, "\n") {
		closeBracket := strings.IndexByte(rawLine, ']')
		if closeBracket < 0 {
			continue
		}
		timestamp, err := timeStamp(rawLine[1:closeBracket])
		if err != nil {
			continue
		}
		text := strings.TrimSpace(rawLine[closeBracket+1:])
		if text != "" {
			lines = append(lines, Line{At: timestamp, Text: text})
		}
	}
	sort.Slice(lines, func(i, j int) bool { return lines[i].At < lines[j].At })
	if len(lines) == 0 {
		return nil, fmt.Errorf("lyrics contain no timed lines")
	}
	return lines, nil
}

func timeStamp(value string) (time.Duration, error) {
	var minutes, seconds, hundredths int
	if _, err := fmt.Sscanf(value, "%d:%d.%d", &minutes, &seconds, &hundredths); err != nil {
		return 0, err
	}
	return time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second +
		time.Duration(hundredths)*10*time.Millisecond, nil
}

func plainLines(raw string) []Line {
	var lines []Line
	for _, text := range strings.Split(raw, "\n") {
		if text = strings.TrimSpace(text); text != "" {
			lines = append(lines, Line{Text: text})
		}
	}
	return lines
}
