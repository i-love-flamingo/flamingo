package web

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"

	"github.com/gorilla/sessions"
)

// Session holds the data connected to the current user session
type Session struct {
	mu              sync.Mutex
	s               *sessions.Session
	hashedid        string
	dirty           map[interface{}]struct{}
	dirtyAll        bool
	sessionSaveMode sessionPersistLevel
}

type sessionPersistLevel uint

const (
	sessionSaveAlways sessionPersistLevel = iota
	sessionSaveOnRead
	sessionSaveOnWrite

	contextSession contextKeyType = "session"
	flashesKey                    = "_flash"
)

// EmptySession creates an empty session instance for testing etc.
func EmptySession() *Session {
	return &Session{s: sessions.NewSession(nil, "")}
}

// ContextWithSession returns a new Context with an attached session
func ContextWithSession(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, contextSession, session)
}

// SessionFromContext allows to retrieve the stored session
func SessionFromContext(ctx context.Context) *Session {
	session, _ := ctx.Value(contextSession).(*Session)
	return session
}

func (s *Session) markDirty(key interface{}) {
	if s.dirty == nil {
		s.dirty = make(map[interface{}]struct{})
	}
	s.dirty[key] = struct{}{}
}

// Load data by a key
func (s *Session) Load(key interface{}) (data interface{}, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, ok = s.s.Values[key]
	if s.sessionSaveMode <= sessionSaveOnRead {
		s.markDirty(key)
	}
	return data, ok
}

// Try to load data by a key
func (s *Session) Try(key interface{}) (data interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, _ = s.s.Values[key]
	if s.sessionSaveMode <= sessionSaveOnRead {
		s.markDirty(key)
	}
	return data
}

// Store data with a key in the Session
func (s *Session) Store(key interface{}, data interface{}) *Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.s.Values[key] = data
	if s.sessionSaveMode <= sessionSaveOnWrite {
		s.markDirty(key)
	}

	return s
}

// Delete a given key from the session
func (s *Session) Delete(key interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.sessionSaveMode <= sessionSaveOnWrite {
		s.markDirty(key)
	}

	delete(s.s.Values, key)
}

// ID returns the Session id
func (s *Session) ID() (id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.s.ID
}

// Keys returns an unordered list of session keys
// Deprecated: please know what you will need
func (s *Session) Keys() []interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	keys := make([]interface{}, len(s.s.Values))
	i := 0
	for k := range s.s.Values {
		keys[i] = k
		i++
	}
	return keys
}

// ClearAll removes all values from the session
// Deprecated: do not use ClearAll
func (s *Session) ClearAll() *Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.s.Values = make(map[interface{}]interface{})
	s.dirtyAll = true
	return s
}

// Flashes returns a slice of flash messages from the session
// todo change?
func (s *Session) Flashes(vars ...string) []interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	// the call to Flashes actually writes to the session
	if s.sessionSaveMode <= sessionSaveOnWrite {
		key := flashesKey
		if len(vars) > 0 {
			key = vars[0]
		}
		s.markDirty(key)
	}

	return s.s.Flashes(vars...)
}

// AddFlash adds a flash message to the session.
// todo change?
func (s *Session) AddFlash(value interface{}, vars ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.sessionSaveMode <= sessionSaveOnWrite {
		key := flashesKey
		if len(vars) > 0 {
			key = vars[0]
		}
		s.markDirty(key)
	}

	s.s.AddFlash(value, vars...)
}

// IDHash - returns the Hashed session id - useful for logs
func (s *Session) IDHash() string {
	if s.hashedid != "" {
		return s.hashedid
	}
	id := s.ID()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hashedid = hashID(id)
	return s.hashedid
}

func hashID(id string) string {
	h := sha256.New()
	h.Write([]byte(id))
	return fmt.Sprintf("%x", h.Sum(nil))
}
