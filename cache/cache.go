package cache

import (
	"container/list"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Cache holds the cache data and a mutex for thread safety.
type Cache[T any] struct {
	data     map[string]*list.Element
	order    *list.List
	mu       sync.Mutex
	capacity int
}

// CacheEntry represents a cached API response entry.
type CacheEntry[T any] struct {
	Key    string
	Record T
	Expiry time.Time
}

// NewCache creates a new Cache instance with a given capacity.
func NewCache[T any](capacity int) *Cache[T] {
	return &Cache[T]{
		data:     make(map[string]*list.Element),
		order:    list.New(),
		capacity: capacity,
	}
}

// RestoreCache loads cache from disk if available.
func RestoreCache[T any](capacity int, cacheDir string) (*Cache[T], error) {
	filename := "cache.json"
	fileLoc := filepath.Join("/app/", cacheDir, filename)

	c := NewCache[T](capacity)

	file, err := os.ReadFile(fileLoc)
	if err != nil {
		// If file doesn't exist, just return empty cache
		if os.IsNotExist(err) {
			return c, nil
		}
		return nil, fmt.Errorf("could not read cache file: %v", err)
	}

	var data map[string]T
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, fmt.Errorf("could not unmarshal cache: %v", err)
	}

	// Populate cache with loaded data (no expiry, or set a default expiry)
	for key, record := range data {
		c.Set(key, record, 24*time.Hour) // Set a default TTL
	}

	return c, nil
}

// Get retrieves a single cached entry if it exists and is not expired.
func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero T
	if elem, found := c.data[key]; found {
		entry := elem.Value.(*CacheEntry[T])
		if time.Now().After(entry.Expiry) {
			// Cache entry expired
			c.order.Remove(elem)
			delete(c.data, key)
			return zero, false
		}
		// Move accessed entry to the front of the list
		c.order.MoveToFront(elem)
		return entry.Record, true
	}
	return zero, false
}

// GetAll returns all non-expired cache entries.
func (c *Cache[T]) GetAll() map[string]T {
	c.mu.Lock()
	defer c.mu.Unlock()

	allEntries := make(map[string]T)
	for key, elem := range c.data {
		entry := elem.Value.(*CacheEntry[T])
		if time.Now().Before(entry.Expiry) {
			allEntries[key] = entry.Record
		}
	}
	return allEntries
}

// BackupCache creates a backup of the cache as a cache.json. file
func (c *Cache[T]) BackupCache(cacheDir string, data map[string]T) error {
	filename := "cache.json"
	fileLoc := filepath.Join("/app/", cacheDir, filename)
	dataBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	if _, err := os.Stat(fileLoc); err == nil {
		if err := os.Remove(fileLoc); err != nil {
			return fmt.Errorf("error removing old cache file: %v", err)
		}
	}

	if err := os.WriteFile(fileLoc, dataBytes, 0644); err != nil {
		return fmt.Errorf("error writing JSON data to cache file: %v", err)
	}
	return nil
}

// Set adds or updates a cached entry with eviction policy.
func (c *Cache[T]) Set(key string, record T, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, found := c.data[key]; found {
		// Update existing entry
		entry := elem.Value.(*CacheEntry[T])
		entry.Record = record
		entry.Expiry = time.Now().Add(ttl)
		c.order.MoveToFront(elem)
		return
	}

	// Evict the oldest entry if the cache is full
	if c.order.Len() >= c.capacity {
		oldest := c.order.Back()
		if oldest != nil {
			oldestEntry := oldest.Value.(*CacheEntry[T])
			c.order.Remove(oldest)
			delete(c.data, oldestEntry.Key)
		}
	}

	// Add new entry
	entry := &CacheEntry[T]{
		Key:    key,
		Record: record,
		Expiry: time.Now().Add(ttl),
	}
	elem := c.order.PushFront(entry)
	c.data[key] = elem
}

// GetTop retrieves the top (most recently used) cached entry if it exists and is not expired.
func (c *Cache[T]) GetTop() (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var zero T
	front := c.order.Front()
	if front == nil {
		return zero, false
	}
	entry := front.Value.(*CacheEntry[T])
	return entry.Record, true
}
