package ytdlp

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestURLCacheDoesNotCompleteReplacedEntry(t *testing.T) {
	originalFetch := fetchDirectURL
	defer func() { fetchDirectURL = originalFetch }()

	firstStarted := make(chan struct{})
	releaseFirst := make(chan struct{})
	var fetchCount atomic.Int32
	fetchDirectURL = func(string) (string, error) {
		call := fetchCount.Add(1)
		if call == 1 {
			close(firstStarted)
			<-releaseFirst
			return "stale-url", nil
		}
		return "fresh-url", nil
	}

	cache := NewURLCache(time.Minute)
	defer cache.Close()

	cache.Prefetch("video")
	select {
	case <-firstStarted:
	case <-time.After(time.Second):
		t.Fatal("first prefetch did not start")
	}

	cache.Invalidate("video")
	cache.Prefetch("video")
	close(releaseFirst)

	url, err := cache.Get("video")
	if err != nil {
		t.Fatalf("Get returned unexpected error: %v", err)
	}
	if url != "fresh-url" {
		t.Fatalf("Get returned %q, want fresh-url", url)
	}
	if calls := fetchCount.Load(); calls < 2 {
		t.Fatalf("fetchDirectURL called %d times, want at least 2", calls)
	}
}

func TestURLCacheGetReturnsClosedError(t *testing.T) {
	cache := NewURLCache(time.Minute)
	cache.Close()

	_, err := cache.Get("video")
	if !errors.Is(err, errURLCacheClosed) {
		t.Fatalf("Get returned error %v, want %v", err, errURLCacheClosed)
	}
}
