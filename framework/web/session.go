package web

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

// todo: still necessary?
const (
	// INFO notice
	INFO = "info"
	// WARNING notice
	WARNING = "warning"
	// ERROR notice
	ERROR = "error"
)

type Session struct {
	mu sync.RWMutex
	//data map[interface{}]interface{}
	//id   string
	s *sessions.Session
}

func newID() string {
	return uuid.New().String()
}

// NewSession wraps a gorilla session
func NewSession(s *sessions.Session) *Session {
	return &Session{s: s}
}

// Load data by a key
func (s *Session) Load(key interface{}) (data interface{}, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, ok = s.s.Values[key]
	return data, ok
}

// Load data by a key
func (s *Session) Try(key interface{}) (data interface{}) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, _ = s.s.Values[key]
	return data
}

// Store data with a key in the Session
func (s *Session) Store(key interface{}, data interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.s.Values[key] = data
}

// Delete a given key from the session
func (s *Session) Delete(key interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.s.Values, key)
}

// ID returns the Session id
func (s *Session) ID() (id string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.s.ID
}

// Migrate should generate a new ID and remove old Session data
func (s *Session) migrate() (id, old_id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	//old_id = s.id
	//s.id = newID()
	return "", ""

	//return s.id, old_id
}

// Flashes returns a slice of flash messages from the session.
//
// A single variadic argument is accepted, and it is optional: it defines
// the flash key. If not defined "_flash" is used by default.
func (s *Session) Flashes(vars ...string) []interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.s.Flashes(vars...)
}

// AddFlash adds a flash message to the session.
//
// A single variadic argument is accepted, and it is optional: it defines
// the flash key. If not defined "_flash" is used by default.
func (s *Session) AddFlash(value interface{}, vars ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.s.AddFlash(value, vars...)
}

// G access gorilla subsession
// deprecated: kept for backwards compatibility
func (s *Session) G() *sessions.Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.s
}
