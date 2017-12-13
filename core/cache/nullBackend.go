package cache

type (
	NullBackend struct{}
)

func (*NullBackend) Get(key string) (entry *CacheEntry, found bool) { return nil, false }
func (*NullBackend) Set(key string, entry *CacheEntry)              {}
func (*NullBackend) Purge(key string)                               {}
func (*NullBackend) PurgeTags(tags []string)                        {}
func (*NullBackend) Flush()                                         {}
