package cache

type (
	NullBackend struct{}
)

func (*NullBackend) Get(key string) (entry *Entry, found bool) { return nil, false }
func (*NullBackend) Set(key string, entry *Entry)              {}
func (*NullBackend) Purge(key string)                          {}
func (*NullBackend) PurgeTags(tags []string)                   {}
func (*NullBackend) Flush()                                    {}
