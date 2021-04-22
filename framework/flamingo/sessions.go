package flamingo

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
	sessionhealthcheck "flamingo.me/flamingo/v3/framework/flamingo/healthcheck"
	"github.com/boj/redistore"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/sessions"
	"github.com/zemirco/memorystore"
)

// SessionModule for session management
type SessionModule struct {
	backend              string
	secret               string
	fileName             string
	secure               bool
	sameSite             string
	storeLength          int
	maxAge               int
	path                 string
	redisHost            string
	redisPassword        string
	redisIdleConnections int
	redisMaxAge          int
	redisDatabase        string
	healthcheckSession   bool
}

// Inject dependencies
func (m *SessionModule) Inject(config *struct {
	// session config is optional to allow usage of the DefaultConfig
	Backend  string `inject:"config:flamingo.session.backend"`
	Secret   string `inject:"config:flamingo.session.secret"`
	FileName string `inject:"config:flamingo.session.file"`
	Secure   bool   `inject:"config:flamingo.session.cookie.secure"`
	SameSite string `inject:"config:flamingo.session.cookie.sameSite"`
	// float64 is used due to the injection as config from json - int is not possible on this
	StoreLength          float64 `inject:"config:flamingo.session.store.length"`
	MaxAge               float64 `inject:"config:flamingo.session.max.age"`
	Path                 string  `inject:"config:flamingo.session.cookie.path"`
	RedisURL             string  `inject:"config:flamingo.session.redis.url"`
	RedisHost            string  `inject:"config:flamingo.session.redis.host"`
	RedisPassword        string  `inject:"config:flamingo.session.redis.password"`
	RedisIdleConnections float64 `inject:"config:flamingo.session.redis.idle.connections"`
	RedisMaxAge          float64 `inject:"config:flamingo.session.redis.maxAge"`
	RedisDatabase        string  `inject:"config:flamingo.session.redis.database,optional"`
	CheckSession         bool    `inject:"config:flamingo.session.healthcheck,optional"`
}) {
	m.backend = config.Backend
	m.secret = config.Secret
	m.fileName = config.FileName
	m.secure = config.Secure
	m.sameSite = config.SameSite
	m.storeLength = int(config.StoreLength)
	m.maxAge = int(config.MaxAge)
	m.path = config.Path
	m.redisHost, m.redisPassword, m.redisDatabase = getRedisConnectionInformation(config.RedisURL, config.RedisHost, config.RedisPassword, config.RedisDatabase)
	m.redisIdleConnections = int(config.RedisIdleConnections)
	m.maxAge = int(config.MaxAge)
	m.healthcheckSession = config.CheckSession
}

// Configure DI
func (m *SessionModule) Configure(injector *dingo.Injector) {
	switch m.backend {
	case "redis":
		var sessionStore *redistore.RediStore
		var err error

		if m.redisDatabase != "" {
			sessionStore, err = redistore.NewRediStoreWithDB(int(m.redisIdleConnections), "tcp", m.redisHost, m.redisPassword, m.redisDatabase, []byte(m.secret))
		} else {
			sessionStore, err = redistore.NewRediStore(int(m.redisIdleConnections), "tcp", m.redisHost, m.redisPassword, []byte(m.secret))
		}

		if err != nil {
			panic(err) // todo: don't panic? fallback?
		}

		sessionStore.SetMaxAge(m.maxAge)
		sessionStore.SetMaxLength(m.storeLength)
		sessionStore.DefaultMaxAge = m.redisMaxAge
		m.setSessionstoreOptions(sessionStore.Options)

		injector.Bind(new(sessions.Store)).ToInstance(sessionStore)
		injector.Bind(new(redis.Pool)).ToInstance(sessionStore.Pool)

		if m.healthcheckSession {
			injector.BindMap(new(healthcheck.Status), "session").To(sessionhealthcheck.RedisSession{})
		}
	case "file":
		os.Mkdir(m.fileName, os.ModePerm)
		sessionStore := sessions.NewFilesystemStore(m.fileName, []byte(m.secret))

		sessionStore.MaxLength(m.storeLength)
		sessionStore.MaxAge(m.maxAge)
		m.setSessionstoreOptions(sessionStore.Options)

		injector.Bind(new(sessions.Store)).ToInstance(sessionStore)

		if m.healthcheckSession {
			injector.BindMap(new(healthcheck.Status), "session").To(sessionhealthcheck.FileSession{})
		}
	default: //memory
		sessionStore := memorystore.NewMemoryStore([]byte(m.secret))

		sessionStore.MaxLength(m.storeLength)
		sessionStore.MaxAge(m.maxAge)
		m.setSessionstoreOptions(sessionStore.Options)

		injector.Bind(new(sessions.Store)).ToInstance(sessionStore)

		if m.healthcheckSession {
			injector.BindMap(new(healthcheck.Status), "session").To(healthcheck.Nil{})
		}
	}
}

func (m *SessionModule) setSessionstoreOptions(options *sessions.Options) {
	options.Domain = ""
	options.Path = m.path
	options.MaxAge = m.maxAge
	options.Secure = m.secure
	options.HttpOnly = true
	switch m.sameSite {
	case "strict":
		options.SameSite = http.SameSiteStrictMode
	case "none":
		options.SameSite = http.SameSiteNoneMode
	case "lax":
		options.SameSite = http.SameSiteLaxMode
	default:
		options.SameSite = http.SameSiteDefaultMode
	}
}

// CueConfig defines the session config scheme
func (*SessionModule) CueConfig() string {
	return `
flamingo: session: {
	backend: *"memory" | "redis" | "file"
	secret: string | *"flamingosecret"
	file: string | *"/sessions"
	store: length: float | int | *(1024 * 1024)
	max: age: float | int | *(60 * 60 * 24 * 30)
	cookie: {
		secure: bool | *true
		path: string | *"/"
		sameSite: *"lax" | "strict" | "none" | "default"
	}
	redis: {
		url: string | *""
		host: string | *"redis"
		password: string | *""
		idle: connections: float | int | *10
		maxAge: float | int | *(60 * 60 * 24 * 30)
		database: string | *""
	}
}
`
}

// FlamingoLegacyConfigAlias maps legacy config to new
func (m *SessionModule) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{
		"session.backend":                "flamingo.session.backend",
		"session.secret":                 "flamingo.session.secret",
		"session.file":                   "flamingo.session.file",
		"session.store.length":           "flamingo.session.store.length",
		"session.max.age":                "flamingo.session.max.age",
		"session.cookie.secure":          "flamingo.session.cookie.secure",
		"session.cookie.path":            "flamingo.session.cookie.path",
		"session.redis.url":              "flamingo.session.redis.url",
		"session.redis.host":             "flamingo.session.redis.host",
		"session.redis.password":         "flamingo.session.redis.password",
		"session.redis.idle.connections": "flamingo.session.redis.idle.connections",
		"session.redis.maxAge":           "flamingo.session.redis.maxAge",
		"core.healthcheck.checkSession":  "flamingo.session.healthcheck",
	}
}

func getRedisConnectionInformation(redisURL, redisHost, redisPassword, redisDatabase string) (string, string, string) {
	if redisURL == "" {
		return redisHost, redisPassword, redisDatabase
	}

	parsedRedisURL, err := url.Parse(redisURL)
	if err != nil {
		return redisHost, redisPassword, redisDatabase
	}

	redisHostFromURL := parsedRedisURL.Host
	if redisHostFromURL != "" {
		redisHost = redisHostFromURL
	}

	redisPasswordFromURL, isRedisPasswordInURL := parsedRedisURL.User.Password()
	if isRedisPasswordInURL {
		redisPassword = redisPasswordFromURL
	}

	redisDatabaseFromPath := strings.Trim(parsedRedisURL.Path, "/")
	redisDatabaseFromQuery := parsedRedisURL.Query().Get("db")
	if len(redisDatabaseFromPath) > 0 {
		redisDatabase = redisDatabaseFromPath
	} else if len(redisDatabaseFromQuery) > 0 {
		redisDatabase = redisDatabaseFromQuery
	}

	return redisHost, redisPassword, redisDatabase
}
