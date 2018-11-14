package store

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/zemirco/memorystore"

	"fmt"

	"flamingo.me/flamingo/core/auth/domain"
)

type (
	Memory struct {
		sync.RWMutex

		store *memorystore.MemoryStore
	}
)

var (
	_ Store = &Memory{}
)

func (s *Memory) Inject(store sessions.Store) {
	if memoryStore, ok := store.(*memorystore.MemoryStore); ok {
		s.store = memoryStore
	} else {
		panic("wrong type provided as memory store")
	}
}

func (s *Memory) DestroySessionsForUser(user domain.User) error {
	s.Lock()
	defer s.Unlock()

	ids, err := s.getAllSessionIds(user)
	if err != nil {
		return err
	}

	s.destroyAllSessionsByIds(ids)
	all := s.store.GetAll()
	delete(all, s.getAllHashesKey(user))

	return nil
}

func (s *Memory) SetHashAndSessionIdForUser(user domain.User, hash string, id string) error {
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
	all := s.store.GetAll()
	all[key] = hash

	return nil
}

func (s *Memory) GetHashByUser(user domain.User) (string, error) {
	key := s.getHashKey(user)
	s.RLock()
	defer s.RUnlock()

	all := s.store.GetAll()
	hash, ok := all[key]
	if !ok {
		return "", fmt.Errorf("there is no hash for user sub %s", user.Sub)
	}

	return hash, nil
}

func (s *Memory) getAllSessionIds(user domain.User) ([]string, error) {
	key := s.getAllHashesKey(user)

	all := s.store.GetAll()
	data, ok := all[key]
	if !ok {
		return []string{}, nil
	}

	var ids []string
	err := json.Unmarshal([]byte(data), &ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *Memory) destroyAllSessionsByIds(ids []string) {
	all := s.store.GetAll()
	for _, id := range ids {
		delete(all, id)
	}
}

func (s *Memory) addSessionsId(user domain.User, ids []string, id string) error {
	ids = append(ids, id)
	data, err := json.Marshal(ids)
	if err != nil {
		return err
	}

	key := s.getAllHashesKey(user)
	all := s.store.GetAll()
	all[key] = string(data)

	return nil
}

func (s *Memory) getHashKey(user domain.User) string {
	return "user_hash_" + user.Sub
}

func (s *Memory) getAllHashesKey(user domain.User) string {
	return "user_hashes_" + user.Sub
}
