package session

import (
	"os"

	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"github.com/boj/redistore"
	"github.com/gorilla/sessions"
	"github.com/zemirco/memorystore"
)

// Module for session management
type Module struct {
	// session config is optional to allow usage of the DefaultConfig
	Backend  string `inject:"config:session.backend"`
	Secret   string `inject:"config:session.secret"`
	FileName string `inject:"config:session.file"`
	Secure   bool   `inject:"config:session.cookie.secure"`
	// float64 is used due to the injection as config from json - int is not possible on this
	StoreLength          float64 `inject:"config:session.store.length"`
	MaxAge               float64 `inject:"config:session.max.age"`
	Path                 string  `inject:"config:session.cookie.path"`
	RedisHost            string  `inject:"config:session.redis.host"`
	RedisPassword        string  `inject:"config:session.redis.password"`
	RedisIdleConnections float64 `inject:"config:session.redis.idle.connections"`
	RedisMaxAge          float64 `inject:"config:session.redis.maxAge"`
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	switch m.Backend {
	case "redis":
		sessionStore, err := redistore.NewRediStore(int(m.RedisIdleConnections), "tcp", m.RedisHost, m.RedisPassword, []byte(m.Secret))
		if err != nil {
			panic(err) // todo: don't panic? fallback?
		}

		sessionStore.SetMaxAge(int(m.MaxAge))
		sessionStore.SetMaxLength(int(m.StoreLength))
		sessionStore.Options.Secure = m.Secure
		sessionStore.Options.HttpOnly = true
		sessionStore.Options.Path = m.Path
		sessionStore.DefaultMaxAge = int(m.RedisMaxAge)

		injector.Bind((*sessions.Store)(nil)).ToInstance(sessionStore)
	case "file":
		os.Mkdir(m.FileName, os.ModePerm)
		sessionStore := sessions.NewFilesystemStore(m.FileName, []byte(m.Secret))

		sessionStore.MaxLength(int(m.StoreLength))
		sessionStore.MaxAge(int(m.MaxAge))
		sessionStore.Options.Secure = m.Secure
		sessionStore.Options.HttpOnly = true
		sessionStore.Options.Path = m.Path

		injector.Bind((*sessions.Store)(nil)).ToInstance(sessionStore)
	default: //memory
		sessionStore := memorystore.NewMemoryStore([]byte(m.Secret))

		sessionStore.MaxLength(int(m.StoreLength))
		sessionStore.MaxAge(int(m.MaxAge))
		sessionStore.Options.Secure = m.Secure
		sessionStore.Options.HttpOnly = true
		sessionStore.Options.Path = m.Path

		injector.Bind((*sessions.Store)(nil)).ToInstance(sessionStore)
	}
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"session.backend":                "memory",
		"session.secret":                 "flamingosecret",
		"session.file":                   "/sessions",
		"session.store.length":           float64(1024 * 1024),
		"session.max.age":                float64(60 * 60 * 24 * 30),
		"session.cookie.secure":          true,
		"session.cookie.path":            "/",
		"session.redis.host":             "redis",
		"session.redis.password":         "",
		"session.redis.idle.connections": float64(10),
		"session.redis.maxAge":           float64(60 * 60 * 24), // 1 day
	}
}
