package store

import (
	"encoding/json"
	"sync"

	"fmt"

	"flamingo.me/flamingo/v3/core/auth/domain"
	"github.com/garyburd/redigo/redis"
)

type (
	Redis struct {
		sync.RWMutex

		pool   *redis.Pool
		maxAge float64
	}
)

var (
	_ Store = &Redis{}
)

func (s *Redis) Inject(pool *redis.Pool, cfg *struct {
	MaxAge float64 `inject:"config:session.max.age"`
}) {
	s.pool = pool
	s.maxAge = cfg.MaxAge
}

func (s *Redis) DestroySessionsForUser(user domain.User) error {
	s.Lock()
	defer s.Unlock()

	ids, err := s.getAllSessionIds(user)
	if err != nil {
		return err
	}

	err = s.destroyAllSessionsByIds(ids)
	if err != nil {
		return err
	}

	conn := s.pool.Get()
	defer conn.Close()
	if _, err := conn.Do("DEL", s.getAllHashesKey(user)); err != nil {
		return err
	}

	return nil
}

func (s *Redis) SetHashAndSessionIdForUser(user domain.User, hash string, id string) error {
	s.Lock()
	defer s.Unlock()

	ids, err := s.getAllSessionIds(user)
	if err != nil {
		return err
	}

	err = s.addSessionsId(user, ids, id)
	if err != nil {
		return err
	}

	key := s.getHashKey(user)
	conn := s.pool.Get()
	defer conn.Close()
	_, err = conn.Do("SETEX", key, int(s.maxAge), hash)

	return nil
}

func (s *Redis) GetHashByUser(user domain.User) (string, error) {
	key := s.getHashKey(user)
	s.RLock()
	defer s.RUnlock()

	conn := s.pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return "", err
	}

	data, err := conn.Do("GET", key)
	if err != nil {
		return "", err
	}

	if data == nil {
		return "", fmt.Errorf("there is no hash for user sub %s", user.Sub)
	}

	b, err := redis.Bytes(data, err)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (s *Redis) getAllSessionIds(user domain.User) ([]string, error) {
	key := s.getAllHashesKey(user)

	conn := s.pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return nil, err
	}

	data, err := conn.Do("GET", key)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return []string{}, nil
	}
	b, err := redis.Bytes(data, err)
	if err != nil {
		return nil, err
	}

	var ids []string
	err = json.Unmarshal(b, &ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *Redis) addSessionsId(user domain.User, ids []string, id string) error {
	ids = append(ids, id)
	data, err := json.Marshal(ids)
	if err != nil {
		return err
	}

	key := s.getAllHashesKey(user)

	maxAge := int(s.maxAge)
	if maxAge == 0 {
		maxAge = 20 * 60
	}

	conn := s.pool.Get()
	defer conn.Close()
	_, err = conn.Do("SETEX", key, maxAge, data)

	return err
}

func (s *Redis) destroyAllSessionsByIds(ids []string) error {
	conn := s.pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return err
	}

	for _, id := range ids {
		if _, err := conn.Do("DEL", "session_"+id); err != nil {
			return err
		}
	}

	return nil
}

func (s *Redis) getHashKey(user domain.User) string {
	return "user_hash_" + user.Sub
}

func (s *Redis) getAllHashesKey(user domain.User) string {
	return "user_hashes_" + user.Sub
}
