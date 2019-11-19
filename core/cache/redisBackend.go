/**
 * @TODO:
 * - write documentation
 */

package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"regexp"
	"runtime"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/gomodule/redigo/redis"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

type (
	// RedisBackend instance representation
	RedisBackend struct {
		BackendMetrics   BackendMetrics
		Pool             *redis.Pool
		logger           flamingo.Logger
		writeLockEnabled bool
	}

	// redisCacheEntryMeta representation
	redisCacheEntryMeta struct {
		Lifetime, Gracetime time.Duration
	}

	// RedisCacheEntry representation
	RedisCacheEntry struct {
		Meta redisCacheEntryMeta
		Data interface{}
	}

	// RedisBackendOptions ...
	RedisBackendOptions struct {
		Network string
		// Redis Host (default 127.0.0.1)
		Host string
		// Redis Port (default 6379)
		Port string
		// Redis Database Number to use (default 0)
		Db int
		// Passwort for the rdis connection. empty string for none (default empty)
		Password string
		// Maximum number of idle redis connections in the pool. (default 8)
		MaxIdle int
		// Timout to close idle connections (default 5m)
		IdleTimeout time.Duration
		// WriteLock feature switch
		WriteLock bool
	}
)

const (
	lockPrefix  = "lock:"
	tagPrefix   = "tag:"
	valuePrefix = "value:"
)

var (
	redisKeyRegex = regexp.MustCompile(`[^a-zA-Z0-9]`)
)

func redisConnector(network, address, password string) (redis.Conn, error) {
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
	return c, err
}

func finalizer(b *RedisBackend) {
	b.close()
}

// DefaultRedisBackendOptions gets the default options for redis backend
func DefaultRedisBackendOptions() RedisBackendOptions {
	return RedisBackendOptions{
		Network:     "tcp",
		Host:        "127.0.0.1",
		Port:        "6379",
		Db:          0,
		Password:    "",
		MaxIdle:     8,
		IdleTimeout: time.Minute * 30,
		WriteLock:   false,
	}
}

// NewRedisBackend creates an redis cache backend
func NewRedisBackend(options RedisBackendOptions, logger flamingo.Logger) *RedisBackend {
	err := mergo.Merge(&options, DefaultRedisBackendOptions())
	if err != nil {
		logger.WithField("category", "redisBackend").Error(fmt.Sprintf("Error merging options: %v", err))
	}

	b := &RedisBackend{
		Pool: &redis.Pool{
			MaxIdle:     options.MaxIdle,
			IdleTimeout: options.IdleTimeout,
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
			Dial: func() (redis.Conn, error) {
				return redisConnector(options.Network, fmt.Sprintf("%v:%v", options.Host, options.Port), options.Password)
			},
		},
		logger:         logger,
		BackendMetrics: NewBackendMetrics("redis"),
	}

	runtime.SetFinalizer(b, finalizer) // close all connections on destruction

	return b
}

// Inject RedisBackend dependencies
func (b *RedisBackend) Inject(logger flamingo.Logger) {
	b.logger = logger
}

// Close ensures all redis connections are closed
func (b *RedisBackend) close() {
	b.Pool.Close()
}

// createPrefixedKey creates an redis-compatible key
func (b *RedisBackend) createPrefixedKey(key string, prefix string) string {
	key = redisKeyRegex.ReplaceAllString(key, "-")
	return fmt.Sprintf("%v%v", prefix, key)
}

// Get an cache key
func (b *RedisBackend) Get(key string) (entry *Entry, found bool) {
	conn := b.Pool.Get()
	defer conn.Close()

	reply, err := conn.Do("GET", b.createPrefixedKey(key, valuePrefix))
	if err != nil {
		b.BackendMetrics.countMiss()
		b.logger.WithField("category", "redisBackend").Error(fmt.Sprintf("Error getting key '%v': %v", key, err))
		return nil, false
	}
	if reply == nil {
		b.BackendMetrics.countMiss()
		b.BackendMetrics.countError("NilReply")
		b.logger.WithField("category", "redisBackend").Error(fmt.Sprintf("Returned nil for key: %v", key))
		return nil, false
	}

	value, err := redis.Bytes(reply, err)
	if err != nil {
		b.BackendMetrics.countError("ByteConvertFailed")
		b.logger.WithField("category", "redisBackend").Error(fmt.Sprintf("Error convert value to bytes of key '%v': %v", key, err))
		return nil, false
	}

	redisEntry, err := b.decodeEntry(value)
	if err != nil {
		b.BackendMetrics.countError("DecodeFailed")
		b.logger.WithField("category", "redisBackend").Error(fmt.Sprintf("Error decoding content of key '%v': %v", key, err))
		return nil, false
	}

	b.BackendMetrics.countHit()
	return b.buildResult(redisEntry), true
}

// Set an cache key
func (b *RedisBackend) Set(key string, entry *Entry) error {
	conn := b.Pool.Get()
	defer conn.Close()

	if b.writeLockEnabled {
		err := b.lock(conn, key)
		if err != nil {
			return err
		}
		defer func() {
			err := b.unLock(conn, key)
			if err != nil {
				b.logger.WithField("category", "redisBackend").Error("Error unlocking key: %v: %v", key, err)
			}
		}()
	}

	redisEntry := b.buildEntry(entry)

	buffer, err := b.encodeEntry(redisEntry)
	if err != nil {
		b.BackendMetrics.countError("EncodeFailed")
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
		b.BackendMetrics.countError("SetFailed")
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
			b.BackendMetrics.countError("SetTagFailed")
			b.logger.WithField("category", "redisBackend").Error("Error setting tag: %v: %v", tag, key)
			return err
		}
	}

	conn.Flush()
	return nil
}

// Purge an cache key
func (b *RedisBackend) Purge(key string) error {
	conn := b.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", b.createPrefixedKey(key, valuePrefix))
	if err != nil {
		return err
	}

	return nil
}

// PurgeTags purges all keys+tags by tag(s)
func (b *RedisBackend) PurgeTags(tags []string) error {
	conn := b.Pool.Get()
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
func (b *RedisBackend) Flush() error {
	conn := b.Pool.Get()
	defer conn.Close()

	err := conn.Send("FLUSHALL")
	if err != nil {
		b.logger.WithField("category", "redisBackend").Error(fmt.Sprintf("Failed purge all keys %v", err))
		return err
	}

	conn.Flush()

	return nil
}

func (b *RedisBackend) encodeEntry(entry *RedisCacheEntry) (*bytes.Buffer, error) {
	gob.Register(entry)
	gob.Register(entry.Data)
	gob.Register(entry.Meta)

	buffer := new(bytes.Buffer)
	err := gob.NewEncoder(buffer).Encode(entry)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (b *RedisBackend) decodeEntry(content []byte) (*RedisCacheEntry, error) {
	buffer := bytes.NewBuffer(content)
	decoder := gob.NewDecoder(buffer)
	entry := new(RedisCacheEntry)
	err := decoder.Decode(&entry)
	if err != nil {
		return nil, err
	}

	return entry, err
}

// buildEntry removes unneeded Meta.Tags before encoding
func (b *RedisBackend) buildEntry(entry *Entry) *RedisCacheEntry {
	return &RedisCacheEntry{
		Meta: redisCacheEntryMeta{
			Lifetime:  entry.Meta.Lifetime,
			Gracetime: entry.Meta.Gracetime,
		},
		Data: entry.Data,
	}
}

// buildResult removes unneeded Meta.Tags before encoding
func (b *RedisBackend) buildResult(entry *RedisCacheEntry) *Entry {
	return &Entry{
		Meta: Meta{
			Lifetime:  entry.Meta.Lifetime,
			Gracetime: entry.Meta.Gracetime,
		},
		Data: entry.Data,
	}
}

// Lock tries to get an remote lock for an key
func (b *RedisBackend) lock(conn redis.Conn, key string) (err error) {
	reply, err := conn.Do(
		"EXISTS",
		b.createPrefixedKey(key, lockPrefix),
	)
	lockValue, err := redis.String(reply, err)
	if "1" == lockValue || err != nil {
		return errors.New(fmt.Sprintf("Lock for key %v already exists", key))
	}

	_, err = conn.Do(
		"SETEX",
		b.createPrefixedKey(key, lockPrefix),
		15,
		"",
	)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get log for key '%v' with error %v", key, err))
	}

	return nil
}

// UnLock releases an Lock for an key
func (b *RedisBackend) unLock(conn redis.Conn, key string) (err error) {
	_, err = conn.Do(
		"DEL",
		b.createPrefixedKey(key, lockPrefix),
	)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to releasing a lock for key '%v' with error %v", key, err))
	}

	return nil
}
