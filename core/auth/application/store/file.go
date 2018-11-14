package store

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	"flamingo.me/flamingo/core/auth/domain"
)

type (
	File struct {
		sync.RWMutex

		path string
	}
)

func (s *File) Inject(cfg *struct {
	Path string `inject:"config:session.file"`
}) {
	s.path = cfg.Path
}

func (s *File) SetHashForUser(user domain.User, hash string) error {
	filename := s.getFileName(user)
	s.Lock()
	defer s.Unlock()

	return ioutil.WriteFile(filename, []byte(hash), 0600)
}

func (s *File) GetHashByUser(user domain.User) (string, error) {
	filename := s.getFileName(user)
	s.RLock()
	defer s.RUnlock()

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *File) getFileName(user domain.User) string {
	return filepath.Join(s.path, "user_"+user.Sub)
}
