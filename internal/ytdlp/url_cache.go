package ytdlp

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
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
	jobs    chan string
	closeCh chan struct{}
	once    sync.Once
}

const prefetchWorkers = 4

func NewURLCache(ttl time.Duration) *URLCache {
	c := &URLCache{
		entries: make(map[string]*urlCacheEntry),
		ttl:     ttl,
		jobs:    make(chan string, 64),
		closeCh: make(chan struct{}),
	}

	// Prefetcher dédié : pool de workers bornés au lieu de
	// goroutines non limitées à chaque appel de Prefetch.
	for i := 0; i < prefetchWorkers; i++ {
		go c.prefetchWorker()
	}

	go c.cleanupLoop()

	return c
}

func (c *URLCache) prefetchWorker() {
	for {
		select {
		case videoID, ok := <-c.jobs:
			if !ok {
				return
			}
			c.mu.Lock()
			entry, exists := c.entries[videoID]
			c.mu.Unlock()
			if !exists {
				continue
			}
			url, err := FetchDirectURL(videoID)
			entry.url = url
			entry.err = err
			close(entry.ready)
		case <-c.closeCh:
			return
		}
	}
}

func (c *URLCache) cleanupLoop() {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()
			for id, entry := range c.entries {
				select {
				case <-entry.ready:
					if now.After(entry.expiry) {
						delete(c.entries, id)
					}
				default:
					// fetch encore en cours, on ne touche pas
				}
			}
			c.mu.Unlock()
		case <-c.closeCh:
			return
		}
	}
}

// Close arrête proprement les workers et le nettoyage périodique.
// À appeler une fois, en fin de vie du programme.
func (c *URLCache) Close() {
	c.once.Do(func() {
		close(c.closeCh)
	})
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

	select {
	case c.jobs <- videoID:
	default:
		go func() {
			url, err := FetchDirectURL(videoID)
			entry.url = url
			entry.err = err
			close(entry.ready)
		}()
	}
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

func (c *URLCache) PrefetchWindow(tracks []Track, centerIdx int) {
	if c == nil || len(tracks) == 0 {
		return
	}
	const ahead = 3
	const behind = 1
	start := centerIdx - behind
	if start < 0 {
		start = 0
	}
	end := centerIdx + ahead
	if end > len(tracks)-1 {
		end = len(tracks) - 1
	}
	for i := start; i <= end; i++ {
		c.Prefetch(tracks[i].ID)
	}
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

// normalizeQuery uniformise la requête avant hachage pour que des
// variations triviales (casse, espaces) partagent la même entrée.
func normalizeQuery(query string) string {
	return strings.ToLower(strings.TrimSpace(query))
}

func (c *SearchCache) getKey(query string) string {
	hash := sha256.Sum256([]byte(normalizeQuery(query)))
	return hex.EncodeToString(hash[:])
}

func (c *SearchCache) Get(query string) ([]Track, bool) {
	key := c.getKey(query)

	c.mu.RLock()
	entry, ok := c.entries[key]
	if !ok {
		c.mu.RUnlock()
		return nil, false
	}
	expired := time.Now().After(entry.Expiry)
	results := entry.Results
	c.mu.RUnlock()

	if expired {
		c.mu.Lock()
		if e, ok := c.entries[key]; ok && time.Now().After(e.Expiry) {
			delete(c.entries, key)
		}
		c.mu.Unlock()
		return nil, false
	}

	return results, true
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
