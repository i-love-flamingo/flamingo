package healthcheck

import (
	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
	"github.com/gomodule/redigo/redis"
)

// RedisSession pool status check
type RedisSession struct {
	pool *redis.Pool
}

var _ healthcheck.Status = &RedisSession{}

// Inject redis pool for session
func (s *RedisSession) Inject(pool *redis.Pool) {
	s.pool = pool
}

// Status checks if the redis server is available
func (s *RedisSession) Status() (bool, string) {
	conn := s.pool.Get()
	_, err := conn.Do("PING")
	if err == nil {
		return true, "success"
	}

	return false, err.Error()
}
