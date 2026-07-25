package ytdlp

import (
	"sync"
	"time"
)

type urlCacheEntry struct {
	url    string
	err    error
	ready  chan struct{}
	expiry time.Time
}

type URLCache struct {
	mu      sync.Mutex
	entries map[string]*urlCacheEntry
	ttl     time.Duration
}

func NewURLCache(ttl time.Duration) *URLCache {
	return &URLCache{
		entries: make(map[string]*urlCacheEntry),
		ttl:     ttl,
	}
}

func (c *URLCache) Prefetch(videoID string) {
	c.mu.Lock()
	if entry, exists := c.entries[videoID]; exists && time.Now().Before(entry.expiry) {
		c.mu.Unlock()
		return
	}

	entry := &urlCacheEntry{
		ready:  make(chan struct{}),
		expiry: time.Now().Add(c.ttl),
	}
	c.entries[videoID] = entry
	c.mu.Unlock()

	go func() {
		url, err := FetchDirectURL(videoID)
		entry.url = url
		entry.err = err
		close(entry.ready)
	}()
}

func (c *URLCache) Get(videoID string) (string, error) {
	c.mu.Lock()
	entry, exists := c.entries[videoID]
	if !exists || !time.Now().Before(entry.expiry) {
		c.mu.Unlock()
		c.Prefetch(videoID)
		c.mu.Lock()
		entry = c.entries[videoID]
	}
	c.mu.Unlock()

	<-entry.ready
	return entry.url, entry.err
}

func (c *URLCache) Invalidate(videoID string) {
	c.mu.Lock()
	delete(c.entries, videoID)
	c.mu.Unlock()
}
