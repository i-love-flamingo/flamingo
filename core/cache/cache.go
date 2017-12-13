package cache

import (
	"time"
)

type (
	CacheEntry struct {
		Tags      []string
		Lifetime  time.Time
		Gracetime time.Time
		Data      interface{}
	}

	Backend interface {
		Get(key string) (entry *CacheEntry, found bool) // Get a cache entry
		Set(key string, entry *CacheEntry)              // Set a cache entry
		//Peek(key string) (entry *CacheEntry, found bool) // Peek for a cache entry, this should not trigger key-updates or weight/priorities to be changed
		Purge(key string)
		PurgeTags(tags []string)
		Flush()
	}
)
