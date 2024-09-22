package bartender

import (
	"container/list"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Cache struct holds the cache data and a mutex for thread safety
type Cache struct {
	data     map[string]*list.Element
	order    *list.List
	mu       sync.Mutex
	capacity int
}

// CacheEntry represents a cached API response entry
type CacheEntry struct {
	Record DrinkResponse
	Expiry time.Time
}

// NewCache creates a new Cache instance with a given capacity
func NewCache(capacity int) *Cache {
	return &Cache{
		data:     make(map[string]*list.Element),
		order:    list.New(),
		capacity: capacity,
	}
}

// TODO
func RestoreCache(capacity int) *Cache {
	// var data DrinkResponse
	filename := "cache.json"
	dirPath, err := os.Getwd()
	if err != nil {
		fmt.Errorf("could not get cwd: %v", err)
	}

	// Need too parse each JSON record and set it in the newly created cache
	filepath := filepath.Join(dirPath, "/bartender", filename)
	c := NewCache(15)

	fmt.Println(filepath, c.capacity)
	return c
}

// Get retrieves a single cached entry if it exists and is not expired
// Currently not used
func (c *Cache) Get(key string) (DrinkResponse, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.data[key]; found {
		entry := elem.Value.(*CacheEntry)
		if time.Now().After(entry.Expiry) {
			// Cache entry expired
			c.order.Remove(elem)
			delete(c.data, key)
			return DrinkResponse{}, false
		}
		// Move accessed entry to the front of the list
		c.order.MoveToFront(elem)
		return entry.Record, true
	}
	return DrinkResponse{}, false
}

// Returns cache contents
// Not sure if cache is ordered
func (c *Cache) GetAll() map[string]DrinkResponse {
	c.mu.Lock()
	defer c.mu.Unlock()

	allEntries := make(map[string]DrinkResponse)
	for e := c.order.Front(); e != nil; e = e.Next() {
		entry := e.Value.(*CacheEntry)
		if time.Now().Before(entry.Expiry) {
			// Include only non-expired entries
			for key, elem := range c.data {
				if elem == e {
					allEntries[key] = entry.Record
					break
				}
			}
		}
	}
	return allEntries
}

// Creates backup of cache in cache.json
func BackupCache(data map[string]DrinkResponse) error {
	dataBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	filename := "cache.json"
	dirPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get cwd: %v", err)
	}

	filepath := filepath.Join(dirPath, "/bartender", filename)

	_, err = os.Stat(filepath)
	if err == nil {
		os.Remove(filepath)
	}

	err = os.WriteFile(filepath, dataBytes, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON data to new persistent cache file: %v", err)
	}
	return nil
}

// Set adds or updates a cached entry with eviction policy
func (c *Cache) Set(key string, record DrinkResponse, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, found := c.data[record.Name]; found {
		// Update existing entry
		entry := elem.Value.(*CacheEntry)
		entry.Record = record
		entry.Expiry = time.Now().Add(ttl)
		c.order.MoveToFront(elem)
		return
	}

	// Evict the oldest entry if the cache is full
	if c.order.Len() >= c.capacity {
		oldest := c.order.Back()
		if oldest != nil {
			c.order.Remove(oldest)
			delete(c.data, oldest.Value.(*CacheEntry).Record.Name)
		}
	}

	// Add new entry
	entry := &CacheEntry{
		Record: record,
		Expiry: time.Now().Add(ttl),
	}
	elem := c.order.PushFront(entry)
	c.data[record.Name] = elem
}
