package store

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	"encoding/json"
	"os"

	"flamingo.me/flamingo/core/auth/domain"
)

type (
	File struct {
		sync.RWMutex

		path string
	}
)

var (
	_ Store = &File{}
)

func (s *File) Inject(cfg *struct {
	Path string `inject:"config:session.file"`
}) {
	s.path = cfg.Path
}

func (s *File) DestroySessionsForUser(user domain.User) error {
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

	err = os.Remove(s.getAllHashesFileName(user))
	if err != nil {
		return err
	}

	return nil
}

func (s *File) SetHashAndSessionIdForUser(user domain.User, hash string, id string) error {
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

	err = s.addSessionsId(user, ids, id)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.getHashFileName(user), []byte(hash), 0600)
}

func (s *File) GetHashByUser(user domain.User) (string, error) {
	filename := s.getHashFileName(user)
	s.RLock()
	defer s.RUnlock()

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *File) getHashFileName(user domain.User) string {
	return filepath.Join(s.path, "user_hash_"+user.Sub)
}

func (s *File) getAllHashesFileName(user domain.User) string {
	return filepath.Join(s.path, "user_hashes_"+user.Sub)
}

func (s *File) getAllSessionIds(user domain.User) ([]string, error) {
	filename := s.getAllHashesFileName(user)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return []string{}, nil
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var ids []string
	err = json.Unmarshal(data, &ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *File) destroyAllSessionsByIds(ids []string) error {
	for _, id := range ids {
		err := os.Remove(filepath.Join(s.path, "session_"+id))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *File) addSessionsId(user domain.User, ids []string, id string) error {
	ids = append(ids, id)
	data, err := json.Marshal(ids)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.getAllHashesFileName(user), data, 0600)
	if err != nil {
		return err
	}

	return nil
}
