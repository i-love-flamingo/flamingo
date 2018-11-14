package application

import (
	"crypto/sha256"
	"fmt"

	"strings"
	"time"

	"flamingo.me/flamingo/core/auth/application/store"
	"flamingo.me/flamingo/core/auth/domain"
	"github.com/gorilla/sessions"
)

const (
	hashKey = "auth.onlyOneDevice.hash"
)

type (
	Synchronizer interface {
		Insert(user *domain.User, session *sessions.Session) error
		IsActive(user domain.User, session *sessions.Session) (bool, error)
	}

	SynchronizerImpl struct {
		store store.Store
	}
)

func (s *SynchronizerImpl) Inject(store store.Store) {
	s.store = store
}

func (s *SynchronizerImpl) Insert(user *domain.User, session *sessions.Session) error {
	concatenated := strings.Join([]string{user.Sub, time.Now().Format(time.RFC3339Nano)}, "|")
	hashBytes := sha256.Sum256([]byte(concatenated))
	hash := fmt.Sprintf("%x", hashBytes)
	session.Values[hashKey] = hash

	if err := s.store.DestroySessionsForUser(*user); err != nil {
		return err
	}

	return s.store.SetHashAndSessionIdForUser(*user, hash, session.ID)
}

func (s *SynchronizerImpl) IsActive(user domain.User, session *sessions.Session) (bool, error) {
	hash, err := s.store.GetHashByUser(user)
	if err != nil {
		return false, err
	}

	return hash == session.Values[hashKey], nil
}
