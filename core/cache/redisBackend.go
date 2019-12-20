package cache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/gomodule/redigo/redis"
)

type (
	// redisBackend instance representation
	redisBackend struct {
		cacheMetrics Metrics
		pool         *redis.Pool
		logger       flamingo.Logger
	}

	RedisBackendFactory struct {
		logger       flamingo.Logger
		frontendName string
		pool         *redis.Pool
		config       *RedisBackendConfig
	}

	RedisBackendConfig struct {
		MaxIdle            int
		IdleTimeOutSeconds int
		Host               string
		Port               string
	}

	// redisCacheEntryMeta representation
	redisCacheEntryMeta struct {
		Lifetime, Gracetime time.Duration
	}

	// redisCacheEntry representation
	redisCacheEntry struct {
		Meta redisCacheEntryMeta
		Data interface{}
	}
)

const (
	tagPrefix   = "tag:"
	valuePrefix = "value:"
)

var (
	redisKeyRegex = regexp.MustCompile(`[^a-zA-Z0-9]`)
)

func init() {
	gob.Register(new(redisCacheEntry))
	gob.Register(new(cachedResponse))
}

func finalizer(b *redisBackend) {
	b.close()
}

// Inject redisBackend dependencies
func (f *RedisBackendFactory) Inject(logger flamingo.Logger) *RedisBackendFactory {
	f.logger = logger
	return f
}

func (f *RedisBackendFactory) Build() (Backend, error) {
	if f.config != nil {
		if f.config.IdleTimeOutSeconds <= 0 {
			return nil, errors.New("IdleTimeOut must be >0")
		}
		if f.config.Host == "" || f.config.Port == "" {
			return nil, errors.New("Host and Port must set")
		}
		f.pool = &redis.Pool{
			MaxIdle:     f.config.MaxIdle,
			IdleTimeout: time.Second * time.Duration(f.config.IdleTimeOutSeconds),
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
			Dial: func() (redis.Conn, error) {
				return f.redisConnector(
					"tcp",
					fmt.Sprintf("%v:%v", f.config.Host, f.config.Port),
					"",
					0,
				)
			},
		}
	}
	b := &redisBackend{
		pool:         f.pool,
		logger:       f.logger,
		cacheMetrics: NewCacheMetrics("redis", f.frontendName),
	}
	runtime.SetFinalizer(b, finalizer) // close all connections on destruction
	return b, nil
}

//SetFrontendName for redis cache metrics
func (f *RedisBackendFactory) SetFrontendName(frontendName string) *RedisBackendFactory {
	f.frontendName = frontendName
	return f
}

//SetConfig for redis
func (f *RedisBackendFactory) SetConfig(config RedisBackendConfig) *RedisBackendFactory {
	f.config = &config
	return f
}

//SetPool directly - use instead of SetConfig if desired
func (f *RedisBackendFactory) SetPool(pool *redis.Pool) *RedisBackendFactory {
	f.pool = pool
	return f
}

func (f *RedisBackendFactory) redisConnector(network, address, password string, db int) (redis.Conn, error) {
	c, err := redis.Dial(network, address)
	if err != nil {
		return nil, err
	}
	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			return nil, err
		}
	}
	if db != 0 {
		if _, err := c.Do("SELECT", db); err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, err
}

// Close ensures all redis connections are closed
func (b *redisBackend) close() {
	b.pool.Close()
}

// createPrefixedKey creates an redis-compatible key
func (b *redisBackend) createPrefixedKey(key string, prefix string) string {
	key = redisKeyRegex.ReplaceAllString(key, "-")
	return fmt.Sprintf("%v%v", prefix, key)
}

// Get an cache key
func (b *redisBackend) Get(key string) (entry *Entry, found bool) {
	conn := b.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("GET", b.createPrefixedKey(key, valuePrefix))
	if err != nil {
		b.cacheMetrics.countMiss()
		b.logger.WithField("category", "redisBackend").Error(fmt.Sprintf("Error getting key '%v': %v", key, err))
		return nil, false
	}
	if reply == nil {
		b.cacheMetrics.countMiss()
		b.cacheMetrics.countError("NilReply")
		b.logger.WithField("category", "redisBackend").Error(fmt.Sprintf("Returned nil for key: %v", key))
		return nil, false
	}

	value, err := redis.Bytes(reply, err)
	if err != nil {
		b.cacheMetrics.countError("ByteConvertFailed")
		b.logger.WithField("category", "redisBackend").Error(fmt.Sprintf("Error convert value to bytes of key '%v': %v", key, err))
		return nil, false
	}

	redisEntry, err := b.decodeEntry(value)
	if err != nil {
		b.cacheMetrics.countError("DecodeFailed")
		b.logger.WithField("category", "redisBackend").Error(fmt.Sprintf("Error decoding content of key '%v': %v", key, err))
		return nil, false
	}

	b.cacheMetrics.countHit()
	return b.buildResult(redisEntry), true
}

// Set an cache key
func (b *redisBackend) Set(key string, entry *Entry) error {
	conn := b.pool.Get()
	defer conn.Close()

	redisEntry := b.buildEntry(entry)

	buffer, err := b.encodeEntry(redisEntry)
	if err != nil {
		b.cacheMetrics.countError("EncodeFailed")
		b.logger.WithField("category", "redisBackend").Error("Error encoding: %v: %v", key, redisEntry)
		return err
	}

	err = conn.Send(
		"SETEX",
		b.createPrefixedKey(key, valuePrefix),
		int(entry.Meta.gracetime.Sub(time.Now().Round(time.Second))),
		buffer,
	)
	if err != nil {
		b.cacheMetrics.countError("SetFailed")
		b.logger.WithField("category", "redisBackend").Error("Error setting key %v with timeout %v and buffer %v", key, int(entry.Meta.Gracetime.Seconds()), buffer)
		return err
	}

	for _, tag := range entry.Meta.Tags {
		err = conn.Send(
			"SADD",
			b.createPrefixedKey(tag, tagPrefix),
			b.createPrefixedKey(key, valuePrefix),
		)
		if err != nil {
			b.cacheMetrics.countError("SetTagFailed")
			b.logger.WithField("category", "redisBackend").Error("Error setting tag: %v: %v", tag, key)
			return err
		}
	}

	conn.Flush()
	return nil
}

// Purge an cache key
func (b *redisBackend) Purge(key string) error {
	conn := b.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", b.createPrefixedKey(key, valuePrefix))
	if err != nil {
		return err
	}

	return nil
}

// PurgeTags purges all keys+tags by tag(s)
func (b *redisBackend) PurgeTags(tags []string) error {
	conn := b.pool.Get()
	defer conn.Close()

	for _, tag := range tags {
		reply, err := conn.Do("SMEMBERS", b.createPrefixedKey(tag, tagPrefix))
		members, err := redis.Strings(reply, err)
		if err != nil {
			b.logger.WithField("category", "redisBackend").Error(fmt.Sprintf("Failed SMEMBERS for tag '%v': %v", tag, err))
		}

		for _, member := range members {
			_, err = conn.Do("DEL", member)
			if err != nil {
				b.logger.WithField("category", "redisBackend").Error(fmt.Sprintf("Failed DEL for key '%v': %v", member, err))
				return err
			}
		}

		_, err = conn.Do("DEL", fmt.Sprintf("%v", tag))
		if err != nil {
			b.logger.WithField("category", "redisBackend").Error(fmt.Sprintf("Failed DEL for key '%v': %v", tag, err))
			return err
		}
	}
	conn.Flush()

	return nil
}

// Flush the whole cache
func (b *redisBackend) Flush() error {
	conn := b.pool.Get()
	defer conn.Close()

	err := conn.Send("FLUSHALL")
	if err != nil {
		b.logger.WithField("category", "redisBackend").Error(fmt.Sprintf("Failed purge all keys %v", err))
		return err
	}

	conn.Flush()

	return nil
}

func (b *redisBackend) encodeEntry(entry *redisCacheEntry) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	err := gob.NewEncoder(buffer).Encode(entry)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (b *redisBackend) decodeEntry(content []byte) (*redisCacheEntry, error) {
	buffer := bytes.NewBuffer(content)
	decoder := gob.NewDecoder(buffer)
	entry := new(redisCacheEntry)
	err := decoder.Decode(&entry)
	if err != nil {
		return nil, err
	}

	return entry, err
}

// buildEntry removes unneeded Meta.Tags before encoding
func (b *redisBackend) buildEntry(entry *Entry) *redisCacheEntry {
	return &redisCacheEntry{
		Meta: redisCacheEntryMeta{
			Lifetime:  entry.Meta.Lifetime,
			Gracetime: entry.Meta.Gracetime,
		},
		Data: entry.Data,
	}
}

// buildResult removes unneeded Meta.Tags before encoding
func (b *redisBackend) buildResult(entry *redisCacheEntry) *Entry {
	return &Entry{
		Meta: Meta{
			Lifetime:  entry.Meta.Lifetime,
			Gracetime: entry.Meta.Gracetime,
		},
		Data: entry.Data,
	}
}
