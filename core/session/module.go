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
	// session config is optional to allow usage of the DefaultConfig
	Backend  string `inject:"config:session.backend,optional"`
	Secret   string `inject:"config:session.secret,optional"`
	FileName string `inject:"config:session.file,optional"`
	// float64 is used due to the injection as config from json - int is not possible on this
	StoreLength      float64 `inject:"config:session.store.length,optional"`
	MaxAge           float64 `inject:"config:session.max.age,optional"`
	RedisHost        string  `inject:"config:session.redis.host,optional"`
	RedisStoreLength float64 `inject:"config:session.redis.store.length,optional"`
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	switch m.Backend {
	case "redis":
		sessionStore, err := redistore.NewRediStore(int(m.RedisStoreLength), "tcp", m.RedisHost, "", []byte(m.Secret))
		if err != nil {
			panic(err)
		}
		sessionStore.SetMaxLength(int(m.StoreLength))
		sessionStore.SetMaxAge(int(m.MaxAge))
		sessionStore.DefaultMaxAge = int(m.MaxAge)
		injector.Bind((*sessions.Store)(nil)).ToInstance(sessionStore)
	case "file":
		os.Mkdir(m.FileName, os.ModePerm)
		sessionStore := sessions.NewFilesystemStore(m.FileName, []byte(m.Secret))
		sessionStore.MaxLength(int(m.StoreLength))
		sessionStore.MaxAge(int(m.MaxAge))
		injector.Bind((*sessions.Store)(nil)).ToInstance(sessionStore)
	default: //memory
		sessionStore := memorystore.NewMemoryStore([]byte(m.Secret))
		sessionStore.MaxLength(int(m.StoreLength))
		sessionStore.MaxAge(int(m.MaxAge))
		injector.Bind((*sessions.Store)(nil)).ToInstance(sessionStore)
	}
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"session.backend":            "memory",
		"session.secret":             "flamingosecret",
		"session.file":               "/sessions",
		"session.store.length":       1024 * 1024,
		"session.max.age":            60 * 60 * 24 * 30,
		"session.redis.host":         "redis",
		"session.redis.store.length": 10,
	}
}
