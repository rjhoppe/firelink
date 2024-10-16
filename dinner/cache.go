package dinner

import (
	"container/list"
	"strconv"
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
	Record CacheJson
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

func (c *Cache) Get(key string) (CacheJson, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.data[key]; found {
		entry := elem.Value.(*CacheEntry)
		if time.Now().After(entry.Expiry) {
			// Cache entry expired
			c.order.Remove(elem)
			delete(c.data, key)
			return CacheJson{}, false
		}
		// Move accessed entry to the front of the list
		c.order.MoveToFront(elem)
		return entry.Record, true
	}
	return CacheJson{}, false
}

// Set adds or updates a cached entry with eviction policy
func (c *Cache) Set(key string, record CacheJson, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, found := c.data[strconv.Itoa(int(record.ID))]; found {
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
			cacheEntry := oldest.Value.(*CacheEntry)
			delete(c.data, cacheEntry.Record.Title)
		}
	}

	// Add new entry
	entry := &CacheEntry{
		Record: record,
		Expiry: time.Now().Add(ttl),
	}

	elem := c.order.PushFront(entry)
	c.data[strconv.Itoa(int(record.ID))] = elem
}

func (c *Cache) GetAll() map[string]CacheJson {
	c.mu.Lock()
	defer c.mu.Unlock()

	allEntries := make(map[string]CacheJson, len(c.data))
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
