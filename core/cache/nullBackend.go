package cache

type (
	// NullBackend does not store anything
	NullBackend struct{}
)

var (
	_ Backend = &NullBackend{}
)

// Get nothing
func (*NullBackend) Get(key string) (entry *Entry, found bool) { return nil, false }

// Set nothing
func (*NullBackend) Set(key string, entry *Entry) {}

// Purge nothing
func (*NullBackend) Purge(key string) {}

// PurgeTags purges nothing
func (*NullBackend) PurgeTags(tags []string) {}

// Flush nothing
func (*NullBackend) Flush() {}
