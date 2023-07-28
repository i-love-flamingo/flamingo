package silentzap

import (
	"sync"

	"go.uber.org/zap/zapcore"
)

type (
	SilentContext struct {
		mu            sync.RWMutex
		storedEntries []storedEntry
		willWrite     bool
	}

	storedEntry struct {
		CheckedLogEntry *zapcore.CheckedEntry
	}
)

func (c *SilentContext) store(entry *zapcore.CheckedEntry) {
	if c == nil {
		return
	}

	go func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.storedEntries = append(
			c.storedEntries,
			storedEntry{
				CheckedLogEntry: entry,
			},
		)
	}()
}

// get returns stored entries and allows direct writing
func (c *SilentContext) get() []storedEntry {
	if c == nil {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	c.willWrite = true

	return c.storedEntries
}

func (c *SilentContext) isWritingAllowed() bool {
	if c == nil {
		return true
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.willWrite
}
