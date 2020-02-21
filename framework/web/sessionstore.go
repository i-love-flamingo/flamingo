package web

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/boj/redistore"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/zemirco/memorystore"
	"go.opencensus.io/trace"
)

// SessionStore handles flamingo's session loading and storing.
// It currently uses gorilla as a backend.
type SessionStore struct {
	sessionStore    sessions.Store
	sessionName     string
	logger          flamingo.Logger
	sessionSaveMode sessionPersistLevel
}

// Inject dependencies.
func (s *SessionStore) Inject(logger flamingo.Logger, cfg *struct {
	SessionStore sessions.Store `inject:",optional"`
	SessionName  string         `inject:"config:flamingo.session.name,optional"`
	SaveMode     string         `inject:"config:flamingo.session.saveMode"`
}) *SessionStore {
	s.sessionStore = cfg.SessionStore
	s.sessionName = cfg.SessionName
	s.logger = logger
	switch cfg.SaveMode {
	case "OnWrite":
		s.sessionSaveMode = sessionSaveOnWrite
	case "OnRead":
		s.sessionSaveMode = sessionSaveOnRead
	default:
		s.sessionSaveMode = sessionSaveAlways
	}

	return s
}

// LoadByRequest loads a Session from an http.Request (it is expected to find the session cookie there)
func (s *SessionStore) LoadByRequest(ctx context.Context, req *http.Request) (*Session, error) {
	// initialize the session
	if s == nil || s.sessionStore == nil {
		return EmptySession(), nil
	}

	var span *trace.Span
	var err error

	_, span = trace.StartSpan(ctx, "flamingo/web/session/load")
	defer span.End()
	gs, err := s.sessionStore.New(req, s.sessionName)

	span.AddAttributes(trace.StringAttribute(flamingo.LogKeySession, hashID(gs.ID)))

	return &Session{s: gs, sessionSaveMode: s.sessionSaveMode}, err
}

// LoadByID loads a Session from a provided session id
func (s *SessionStore) LoadByID(ctx context.Context, id string) (*Session, error) {
	return s.LoadByRequest(ctx, s.requestFromID(id))
}

func (s *SessionStore) requestFromID(id string) *http.Request {
	var codecs []securecookie.Codec

	switch s := s.sessionStore.(type) {
	case *memorystore.MemoryStore:
		codecs = s.Codecs
	case *redistore.RediStore:
		codecs = s.Codecs
	case *sessions.FilesystemStore:
		codecs = s.Codecs
	default:
		panic("not supported")
	}

	cookie, err := securecookie.EncodeMulti(s.sessionName, id, codecs...)
	if err != nil {
		panic(err)
	}

	return &http.Request{
		Header: map[string][]string{"Cookie": {
			fmt.Sprintf("%s=%s", s.sessionName, cookie),
		}},
	}
}

type headerResponseWriter http.Header

// emptyResponseWriter to be able to properly persist sessions
func (w headerResponseWriter) Header() http.Header {
	return http.Header(w)
}
func (headerResponseWriter) Write([]byte) (int, error)  { return 0, io.ErrUnexpectedEOF }
func (headerResponseWriter) WriteHeader(statusCode int) {}

// Save stores a session back in the session storage.
// The returned headers should be applied, they usually contain SetCookie headers.
func (s *SessionStore) Save(ctx context.Context, session *Session) (http.Header, error) {
	// ensure that the session has been saved in the backend
	if s == nil || s.sessionStore == nil {
		return nil, nil
	}

	session.mu.Lock()
	defer session.mu.Unlock()

	gs := session.s

	// copy dirty values to new instance and move Values to original session
	if s.sessionSaveMode != sessionSaveAlways && !session.dirtyAll && session.s.ID != "" {
		// no dirty data means we do not need to persist anything at all
		if len(session.dirty) == 0 {
			return nil, nil
		}

		newGs, err := s.LoadByID(ctx, session.s.ID)
		if err != nil {
			return nil, err
		}
		var ok bool
		for k := range session.dirty {
			newGs.s.Values[k], ok = gs.Values[k]
			if !ok {
				delete(newGs.s.Values, k)
			}
		}
		gs.Values = newGs.s.Values
		session.dirty = nil
	}

	_, span := trace.StartSpan(ctx, "flamingo/web/session/save")
	defer span.End()
	rw := headerResponseWriter(make(http.Header))
	if err := s.sessionStore.Save(s.requestFromID(session.s.ID), rw, gs); err != nil {
		return nil, err
	}

	return rw.Header(), nil
}

// AddHTTPHeader adds the sources http.Header to the target.
func AddHTTPHeader(target, source http.Header) {
	for k, v := range source {
		for _, v := range v {
			target.Add(k, v)
		}
	}
}
