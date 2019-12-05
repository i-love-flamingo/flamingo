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
)

type (
	// RedisBackend instance representation
	RedisBackend struct {
		cacheMetrics CacheMetrics
		pool         *redis.Pool
		logger       flamingo.Logger
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
)

const (
	tagPrefix   = "tag:"
	valuePrefix = "value:"
)

var (
	redisKeyRegex = regexp.MustCompile(`[^a-zA-Z0-9]`)
)

func init() {
	gob.Register(new(RedisCacheEntry))
}

func finalizer(b *RedisBackend) {
	b.close()
}

// NewRedisBackend creates an redis cache backend
func NewRedisBackend(pool *redis.Pool, frontendName string) *RedisBackend {
	b := &RedisBackend{
		pool:         pool,
		logger:       flamingo.NullLogger{},
		cacheMetrics: NewCacheMetrics("redis", frontendName),
	}

	runtime.SetFinalizer(b, finalizer) // close all connections on destruction

	return b
}

// Inject RedisBackend dependencies
func (b *RedisBackend) Inject(pool *redis.Pool, frontendName string, logger flamingo.Logger) {
	b.pool = pool
	b.cacheMetrics = NewCacheMetrics("redis", frontendName)
	b.logger = logger
}

// Close ensures all redis connections are closed
func (b *RedisBackend) close() {
	b.pool.Close()
}

// createPrefixedKey creates an redis-compatible key
func (b *RedisBackend) createPrefixedKey(key string, prefix string) string {
	key = redisKeyRegex.ReplaceAllString(key, "-")
	return fmt.Sprintf("%v%v", prefix, key)
}

// Get an cache key
func (b *RedisBackend) Get(key string) (entry *Entry, found bool) {
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
func (b *RedisBackend) Set(key string, entry *Entry) error {
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
func (b *RedisBackend) Purge(key string) error {
	conn := b.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", b.createPrefixedKey(key, valuePrefix))
	if err != nil {
		return err
	}

	return nil
}

// PurgeTags purges all keys+tags by tag(s)
func (b *RedisBackend) PurgeTags(tags []string) error {
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
func (b *RedisBackend) Flush() error {
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

func (b *RedisBackend) encodeEntry(entry *RedisCacheEntry) (*bytes.Buffer, error) {
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
