package healthcheck

import (
	"github.com/garyburd/redigo/redis"
)

type (
	RedisSession struct {
		pool *redis.Pool
	}
)

var (
	_ Status = &RedisSession{}
)

func (s *RedisSession) Inject(pool *redis.Pool) {
	s.pool = pool
}

func (s *RedisSession) Status() (bool, string) {
	conn := s.pool.Get()
	_, err := conn.Do("PING")
	if err == nil {
		return true, "success"
	}

	return false, err.Error()
}
