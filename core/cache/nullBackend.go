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
func (*NullBackend) Set(key string, entry *Entry) error { return nil }

// Purge nothing
func (*NullBackend) Purge(key string) error { return nil }

// PurgeTags purges nothing
func (*NullBackend) PurgeTags(tags []string) error { return nil }

// Flush nothing
func (*NullBackend) Flush() error { return nil }

// FlushSupport returns false, because the Backend doesn't support
func (*NullBackend) FlushSupport() bool { return false }

// TagSupport returns false, because the Backend doesn't support
func (*NullBackend) TagSupport() bool { return false }
