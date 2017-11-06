package session

import (
	"os"

	"github.com/boj/redistore"
	"github.com/gorilla/sessions"
	"github.com/zemirco/memorystore"
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
)

// Module for session management
type Module struct {
	Backend          string `inject:"config:session.backend"`
	Secret           string `inject:"config:session.secret"`
	RedisHost        string `inject:"config:session.redis.host"`
	FileName         string `inject:"config:session.file"`
	StoreLength      int    `inject:"config:session.store.length"`
	RedisStoreLength int    `inject:"config:session.redis.store.length"`
	MaxAge           int    `inject:"config:session.redis.max.age"`
}

func (m *Module) Configure(injector *dingo.Injector) {
	switch m.Backend {
	case "redis":
		sessionStore, err := redistore.NewRediStore(m.RedisStoreLength, "tcp", m.RedisHost, "", []byte(m.Secret))
		if err != nil {
			panic(err)
		}
		sessionStore.SetMaxLength(m.StoreLength)
		sessionStore.DefaultMaxAge = m.MaxAge
		injector.Bind((*sessions.Store)(nil)).ToInstance(sessionStore)
	case "file":
		os.Mkdir(m.FileName, os.ModePerm)
		sessionStore := sessions.NewFilesystemStore(m.FileName, []byte(m.Secret))
		sessionStore.MaxLength(m.StoreLength)
		injector.Bind((*sessions.Store)(nil)).ToInstance(sessionStore)
	default: //memory
		sessionStore := memorystore.NewMemoryStore([]byte(m.Secret))
		sessionStore.MaxLength(m.StoreLength)
		injector.Bind((*sessions.Store)(nil)).ToInstance(sessionStore)
	}
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"session.backend":            "memory",
		"session.secret":             "flamingosecret",
		"session.redis.host":         "redis",
		"session.file":               "/sessions",
		"session.store.length":       1024 * 1024,
		"session.redis.store.length": 10,
		"session.redis.max.age":      60 * 60 * 24 * 30,
	}
}
