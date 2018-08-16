package cache

import (
	"time"
)

type (
	// Meta describes life and gracetimes, as well as tags, for cache entries
	Meta struct {
		Tags                []string
		Lifetime, Gracetime time.Duration
		lifetime, gracetime time.Time
	}

	// Entry is a cached object with associated meta data
	Entry struct {
		Meta Meta
		Data interface{}
	}

	// Backend describes a cache backend, responsible for storing, flushing, setting and getting entries
	Backend interface {
		Get(key string) (entry *Entry, found bool) // Get a cache entry
		Set(key string, entry *Entry) error        // Set a cache entry
		//Peek(key string) (entry *CacheEntry, found bool) // Peek for a cache entry, this should not trigger key-updates or weight/priorities to be changed
		Purge(key string) error
		PurgeTags(tags []string) error
		Flush() error
	}

	loaderResponse struct {
		data interface{}
		meta *Meta
	}
)
