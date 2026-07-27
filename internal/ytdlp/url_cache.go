package ytdlp

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

// ============================================
// URL CACHE
// ============================================

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

func (c *URLCache) Clear() {
	c.mu.Lock()
	c.entries = make(map[string]*urlCacheEntry)
	c.mu.Unlock()
}

func (c *URLCache) Stats() (int, int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	total := len(c.entries)
	active := 0
	for _, entry := range c.entries {
		if time.Now().Before(entry.expiry) {
			active++
		}
	}
	return total, active
}

// ============================================
// SEARCH CACHE
// ============================================

type searchEntry struct {
	Results []Track
	Expiry  time.Time
}

type SearchCache struct {
	mu      sync.RWMutex
	entries map[string]*searchEntry
	ttl     time.Duration
}

var (
	globalSearchCache *SearchCache
	cacheInitOnce     sync.Once
)

func getSearchCache() *SearchCache {
	cacheInitOnce.Do(func() {
		globalSearchCache = &SearchCache{
			entries: make(map[string]*searchEntry),
			ttl:     5 * time.Minute,
		}
	})
	return globalSearchCache
}

func (c *SearchCache) getKey(query string) string {
	hash := sha256.Sum256([]byte(query))
	return hex.EncodeToString(hash[:])
}

func (c *SearchCache) Get(query string) ([]Track, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := c.getKey(query)
	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}

	if time.Now().After(entry.Expiry) {
		return nil, false
	}

	return entry.Results, true
}

func (c *SearchCache) Set(query string, results []Track) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.getKey(query)
	c.entries[key] = &searchEntry{
		Results: results,
		Expiry:  time.Now().Add(c.ttl),
	}
}

func (c *SearchCache) Invalidate(query string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.getKey(query)
	delete(c.entries, key)
}

func (c *SearchCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]*searchEntry)
}

func (c *SearchCache) Stats() (int, int) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := len(c.entries)
	active := 0
	for _, entry := range c.entries {
		if time.Now().Before(entry.Expiry) {
			active++
		}
	}
	return total, active
}

// ============================================
// FONCTIONS GLOBALES POUR LE SEARCH CACHE
// ============================================

func GetSearchCache() *SearchCache {
	return getSearchCache()
}

func ClearSearchCache() {
	getSearchCache().Clear()
}

func GetSearchCacheStats() (int, int) {
	return getSearchCache().Stats()
}