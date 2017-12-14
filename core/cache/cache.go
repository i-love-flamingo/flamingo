package cache

import (
	"time"
)

type (
	Meta struct {
		Tags                []string
		Lifetime, Gracetime time.Duration
		lifetime, gracetime time.Time
	}

	Entry struct {
		Meta Meta
		Data interface{}
	}

	Backend interface {
		Get(key string) (entry *Entry, found bool) // Get a cache entry
		Set(key string, entry *Entry)              // Set a cache entry
		//Peek(key string) (entry *CacheEntry, found bool) // Peek for a cache entry, this should not trigger key-updates or weight/priorities to be changed
		Purge(key string)
		PurgeTags(tags []string)
		Flush()
	}

	loaderResponse struct {
		data interface{}
		meta *Meta
	}
)
