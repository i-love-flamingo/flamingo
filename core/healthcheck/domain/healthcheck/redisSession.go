package healthcheck

import (
	"github.com/garyburd/redigo/redis"
)

type (
	// RedisSession is the healthcheck for the redis session
	RedisSession struct {
		pool *redis.Pool
	}
)

var (
	_ Status = &RedisSession{}
)

// Inject dependencies
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
