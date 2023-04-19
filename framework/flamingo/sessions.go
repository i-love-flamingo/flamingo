package flamingo

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
	sessionhealthcheck "flamingo.me/flamingo/v3/framework/flamingo/healthcheck"
	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
	"github.com/zemirco/memorystore"
)

// SessionModule for session management
type SessionModule struct {
	backend              string
	secret               string
	fileName             string
	secure               bool
	sameSite             http.SameSite
	storeLength          int
	maxAge               int
	path                 string
	redisHost            string
	redisPassword        string
	redisIdleConnections int
	redisDatabase        int
	redisTLS             bool
	redisClusterMode     bool
	redisTimeout         time.Duration
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
	RedisDatabase        int     `inject:"config:flamingo.session.redis.database,optional"`
	RedisTLS             bool    `inject:"config:flamingo.session.redis.tls,optional"`
	RedisClusterMode     bool    `inject:"config:flamingo.session.redis.clusterMode,optional"`
	RedisTimeout         string  `inject:"config:flamingo.session.redis.timeout,optional"`
	CheckSession         bool    `inject:"config:flamingo.session.healthcheck,optional"`
}) {
	m.backend = config.Backend
	m.secret = config.Secret
	m.fileName = config.FileName
	m.secure = config.Secure
	m.storeLength = int(config.StoreLength)
	m.maxAge = int(config.MaxAge)
	m.path = config.Path
	m.redisHost, m.redisPassword, m.redisDatabase = getRedisConnectionInformation(config.RedisURL, config.RedisHost, config.RedisPassword, config.RedisDatabase)
	m.redisIdleConnections = int(config.RedisIdleConnections)
	m.redisTLS = config.RedisTLS
	m.redisClusterMode = config.RedisClusterMode
	m.healthcheckSession = config.CheckSession

	if config.RedisTimeout != "" {
		redisTimeout, err := time.ParseDuration(config.RedisTimeout)
		if err != nil {
			panic(fmt.Errorf("invalid duration on %q: %q (%w)", "flamingo.session.redis.timeout", config.RedisTimeout, err))
		}
		m.redisTimeout = redisTimeout
	}

	switch config.SameSite {
	case "strict":
		m.sameSite = http.SameSiteStrictMode
	case "none":
		m.sameSite = http.SameSiteNoneMode
	case "lax":
		m.sameSite = http.SameSiteLaxMode
	default:
		m.sameSite = http.SameSiteDefaultMode
	}
}

// Configure DI
func (m *SessionModule) Configure(injector *dingo.Injector) {
	switch m.backend {
	case "redis":
		var client redis.UniversalClient

		var tlsConfig *tls.Config
		if m.redisTLS {
			tlsConfig = &tls.Config{}
		}

		if m.redisClusterMode {
			client = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:     []string{m.redisHost},
				Password:  m.redisPassword,
				PoolSize:  m.redisIdleConnections,
				TLSConfig: tlsConfig,
			})
		} else {
			client = redis.NewClient(&redis.Options{
				Addr:      m.redisHost,
				Password:  m.redisPassword,
				DB:        m.redisDatabase,
				PoolSize:  m.redisIdleConnections,
				TLSConfig: tlsConfig,
			})
		}

		ctx := context.Background()
		if m.redisTimeout > 0 {
			c, cancel := context.WithTimeout(ctx, m.redisTimeout)
			defer cancel()
			ctx = c
		}
		sessionStore, err := redisstore.NewRedisStore(ctx, client)
		if err != nil {
			panic(fmt.Errorf("failed on creating redis store: %w", err))
		}

		sessionStore.Options(sessions.Options{
			Path:     m.path,
			MaxAge:   m.maxAge,
			Secure:   m.secure,
			HttpOnly: true,
			SameSite: m.sameSite,
		})
		sessionStore.Serializer(maxLengthSerializer{
			maxLength:  m.storeLength,
			serializer: redisstore.GobSerializer{},
		})

		injector.Bind(new(sessions.Store)).ToInstance(sessionStore)

		injector.Bind(new(redis.UniversalClient)).ToInstance(client)
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
	default: // memory
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
	options.SameSite = m.sameSite
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
		database: float | int | *0
		tls: bool | *false
		clusterMode: bool | *false
		timeout: string | *"5s"
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
		"core.healthcheck.checkSession":  "flamingo.session.healthcheck",
	}
}

func getRedisConnectionInformation(redisURL, redisHost, redisPassword string, redisDatabase int) (string, string, int) {
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
		redisDatabase, err = strconv.Atoi(redisDatabaseFromPath)
		if err != nil {
			panic(err)
		}
	} else if len(redisDatabaseFromQuery) > 0 {
		redisDatabase, err = strconv.Atoi(redisDatabaseFromQuery)
		if err != nil {
			panic(err)
		}
	}

	return redisHost, redisPassword, redisDatabase
}

type maxLengthSerializer struct {
	maxLength  int
	serializer redisstore.SessionSerializer
}

func (m maxLengthSerializer) Serialize(s *sessions.Session) ([]byte, error) {
	b, err := m.serializer.Serialize(s)
	if err != nil {
		return nil, err
	}

	if m.maxLength != 0 && len(b) > m.maxLength {
		return nil, errors.New("the value to store is too big")
	}

	return b, nil
}

func (m maxLengthSerializer) Deserialize(b []byte, s *sessions.Session) error {
	return m.serializer.Deserialize(b, s)
}
