package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	cache := NewCache[string](10)
	cache.Set("test", "test", 10*time.Second)
	val, found := cache.Get("test")
	assert.True(t, found)
	assert.Equal(t, "test", val)
}

func TestCache_GetAll(t *testing.T) {
	cache := NewCache[string](10)
	cache.Set("test1", "test1", 10*time.Second)
	cache.Set("test2", "test2", 10*time.Second)
	cache.Set("test3", "test3", 10*time.Second)
	cache.Set("test4", "test4", 10*time.Second)
	cache.Set("test5", "test5", 10*time.Second)
	all := cache.GetAll()
	assert.Equal(t, 5, len(all))
	assert.Contains(t, all, "test1")
	assert.Contains(t, all, "test2")
	assert.Contains(t, all, "test3")
	assert.Contains(t, all, "test4")
	assert.Contains(t, all, "test5")
}

func TestCache_GetTop(t *testing.T) {
	cache := NewCache[string](10)
	cache.Set("test1", "test1", 10*time.Second)
	cache.Set("test2", "test2", 10*time.Second)
	top, found := cache.GetTop()
	assert.True(t, found)
	assert.Equal(t, "test2", top)
	cache.Set("test3", "test3", 10*time.Second)
	cache.Set("test4", "test4", 10*time.Second)
	cache.Set("test5", "test5", 10*time.Second)
	secondTop, secondFound := cache.GetTop()
	assert.True(t, secondFound)
	assert.Equal(t, "test5", secondTop)
}

func TestCache_Clear(t *testing.T) {
	cache := NewCache[string](10)
	cache.Set("test1", "test1", 10*time.Second)
	cache.Set("test2", "test2", 10*time.Second)
	cache.Clear()
	assert.Equal(t, 0, len(cache.GetAll()))
}
