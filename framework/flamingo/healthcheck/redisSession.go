package healthcheck

import (
	"context"

	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
	"github.com/go-redis/redis/v8"
)

// RedisSession status check
type RedisSession struct {
	client redis.UniversalClient
}

var _ healthcheck.Status = &RedisSession{}

// Inject redis client for session
func (s *RedisSession) Inject(client redis.UniversalClient) {
	s.client = client
}

// Status checks if the redis server is available
func (s *RedisSession) Status() (bool, string) {
	err := s.client.Ping(context.Background()).Err()
	if err != nil {
		return false, err.Error()
	}

	return true, "success"
}
