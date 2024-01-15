package caching

import (
	"encoding/json"
	"time"

	"github.com/maypok86/otter"
)

type AppCacher interface {
	Set(key string, value interface{}) error
	SetWithTTL(key string, value interface{}, duration time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	Clear() error
	Has(key string) (bool, error)
}

type AppCache struct {
	cache *otter.Cache[string, interface{}]
}

func NewAppCache() (*AppCache, error) {
	// Initialize BigCache with desired configuration
	builder, err := otter.NewBuilder[string, interface{}](1000)
	if err != nil {
		panic(err)
	}

	// StatsEnabled determines whether statistics should be calculated when the cache is running.
	// By default, statistics calculating is disabled.
	builder.StatsEnabled(false)

	// Build creates a new cache object or
	// returns an error if invalid parameters were passed to the builder.
	cache, err := builder.Build()
	if err != nil {
		panic(err)
	}

	return &AppCache{cache: cache}, nil
}

func (c *AppCache) Set(key string, value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		// value is not already a byte slice, convert it
		var err error
		bytes, err = json.Marshal(value)
		if err != nil {
			return err
		}
	}
	c.cache.Set(key, bytes)
	return nil
}
func (c *AppCache) SetWithTTL(key string, value interface{}, duration time.Duration) error {
	bytes, ok := value.([]byte)
	if !ok {
		// value is not already a byte slice, convert it
		var err error
		bytes, err = json.Marshal(value)
		if err != nil {
			return err
		}
	}
	c.cache.SetWithTTL(key, bytes, duration)
	return nil
}

func (c *AppCache) Get(key string) (interface{}, error) {
	// Get value from BigCache
	entry, ok := c.cache.Get(key)
	if ok {
		return entry, nil
	}
	return nil, nil
}

func (c *AppCache) Delete(key string) error {

	c.cache.Delete(key)
	return nil
}
func (c *AppCache) DeleteAll(key []string) error {
	for _, k := range key {
		c.cache.Delete(k)
	}
	return nil
}
func (c *AppCache) Clear() error {
	// Clear the entire cache
	c.cache.Clear()
	return nil
}

func (c *AppCache) Has(key string) (bool, error) {
	// Check if key exists in cache
	return c.cache.Has(key), nil
}
