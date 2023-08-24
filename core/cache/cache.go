package cache

import (
	"time"

	"go.opencensus.io/trace"
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
		Get(key string) (entry *Entry, found bool)
		Set(key string, entry *Entry) error
		Purge(key string) error
		PurgeTags(tags []string) error
		Flush() error
	}

	loaderResponse struct {
		data interface{}
		meta *Meta
		span trace.SpanContext
	}
)
