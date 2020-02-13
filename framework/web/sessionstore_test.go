package web

import (
	"context"
	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/stretchr/testify/assert"
	"github.com/zemirco/memorystore"
)

func testsession(t *testing.T, store *SessionStore) (session *Session, saveSession func()) {
	session, err := store.LoadByID(context.Background(), "test-id")
	assert.NoError(t, err)

	return session, func() {
		_, err = store.Save(context.Background(), session)
		assert.NoError(t, err)
	}
}

func TestSessionSaveAlways(t *testing.T) {
	store := memorystore.NewMemoryStore([]byte("flamingosecret"))
	sessionStore := &SessionStore{logger: new(flamingo.StdLogger), sessionName: "test", sessionStore: store}

	// prepare
	session, saveSession := testsession(t, sessionStore)
	session.Store("key1", "val0")
	session.Store("key2", "val0")
	session.Store("key3", "val0")
	saveSession()

	session1, saveSession1 := testsession(t, sessionStore)
	session1.Store("key1", "val1")
	session1.Store("key2", "val1")
	session1.Delete("key3")

	session2, saveSession2 := testsession(t, sessionStore)
	session2.Store("key1", "val2")

	saveSession1()
	saveSession2()

	session, _ = testsession(t, sessionStore)
	assert.Equal(t, map[interface{}]interface{}{
		"key1": "val2", // from session2
		"key2": "val0", // not stored from 1 because 2 overwrote it
		"key3": "val0", // untouched
	}, session.s.Values)
}

func TestSessionSaveOnRead(t *testing.T) {
	store := memorystore.NewMemoryStore([]byte("flamingosecret"))
	sessionStore := &SessionStore{logger: new(flamingo.StdLogger), sessionName: "test", sessionStore: store, sessionSaveMode: sessionSaveOnRead}

	// prepare
	session, saveSession := testsession(t, sessionStore)
	session.Store("key1", "val0")
	session.Store("key2", "val0")
	session.Store("key3", "val0")
	session.Store("key4", "val0")
	saveSession()

	session1, saveSession1 := testsession(t, sessionStore)
	session1.Store("key1", "val1")
	session1.Store("key2", "val1")
	session1.Delete("key3")
	session1.Store("key4", "val1")

	session2, saveSession2 := testsession(t, sessionStore)
	session2.Load("key1")
	session1.Store("key4", "val2")

	saveSession1()
	saveSession2()

	session, _ = testsession(t, sessionStore)
	assert.Equal(t, map[interface{}]interface{}{
		"key1": "val0", // from session2
		"key2": "val1", // from session1
		//"key3": nil,    // deleted
		"key4": "val2", // from session2
	}, session.s.Values)
}

func TestSessionSaveOnWrite(t *testing.T) {
	store := memorystore.NewMemoryStore([]byte("flamingosecret"))
	sessionStore := &SessionStore{logger: new(flamingo.StdLogger), sessionName: "test", sessionStore: store, sessionSaveMode: sessionSaveOnWrite}

	// prepare
	session, saveSession := testsession(t, sessionStore)
	session.Store("key1", "val0")
	session.Store("key2", "val0")
	session.Store("key3", "val0")
	saveSession()

	session1, saveSession1 := testsession(t, sessionStore)
	session1.Store("key1", "val1")
	session1.Store("key2", "val1")
	session1.Delete("key3")

	session2, saveSession2 := testsession(t, sessionStore)
	session2.Store("key1", "val2")
	session2.Load("key2")

	saveSession1()
	saveSession2()

	session, _ = testsession(t, sessionStore)
	assert.Equal(t, map[interface{}]interface{}{
		"key1": "val2", // from session2
		"key2": "val1", // from session1
		//"key3": nil,    // deleted
	}, session.s.Values)
}

func TestSessionSavOnWriteDirtyAll(t *testing.T) {
	store := memorystore.NewMemoryStore([]byte("flamingosecret"))
	sessionStore := &SessionStore{logger: new(flamingo.StdLogger), sessionName: "test", sessionStore: store, sessionSaveMode: sessionSaveOnWrite}

	// prepare
	session, saveSession := testsession(t, sessionStore)
	session.Store("key1", "val0")
	session.Store("key2", "val0")
	session.Store("key3", "val0")
	saveSession()

	session1, saveSession1 := testsession(t, sessionStore)
	session1.Store("key1", "val1")
	session1.Store("key2", "val1")
	session1.Delete("key3")

	session2, saveSession2 := testsession(t, sessionStore)
	session2.ClearAll()

	saveSession1()
	saveSession2()

	session, _ = testsession(t, sessionStore)
	assert.Equal(t, map[interface{}]interface{}{}, session.s.Values)
}
