package session_test

import (
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/zemirco/memorystore"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

// countingStore counts the loads and saves.
type countingStore struct {
	sessions.Store

	loads atomic.Int32
	saves atomic.Int32
}

func (s *countingStore) New(req *http.Request, name string) (*sessions.Session, error) {
	s.loads.Add(1)

	return s.Store.New(req, name)
}

func (s *countingStore) Save(req *http.Request, rw http.ResponseWriter, session *sessions.Session) error {
	s.saves.Add(1)

	return s.Store.Save(req, rw, session)
}

// countingBackend is a flamingo.CustomSessionBackend with a countingStore.
type countingBackend struct {
	store *countingStore
}

func (b countingBackend) SessionStore() sessions.Store {
	return b.store
}

func TestBackend_CustomBackendUsed(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	store := &countingStore{Store: memorystore.NewMemoryStore([]byte("secret"))}

	module := dingo.ModuleFunc(func(injector *dingo.Injector) {
		injector.Bind(new(flamingo.CustomSessionBackend)).ToInstance(countingBackend{store: store})
	})

	app := bootSession(t, []dingo.Module{module}, "flamingo.session.backend: custom")

	app.Get(t, "/en/basket/count")

	assert.Positive(t, store.loads.Load(), "loads")
	assert.Positive(t, store.saves.Load(), "saves")
}
