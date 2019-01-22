package application

import (
	"crypto/sha256"
	"fmt"

	"strings"
	"time"

	"flamingo.me/flamingo/v3/core/auth/application/store"
	"flamingo.me/flamingo/v3/core/auth/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/gorilla/sessions"
)

const (
	hashKey = "auth.preventSimultaneousSessions.hash"
)

type (
	Synchronizer interface {
		Insert(user domain.User, session *sessions.Session) error
		IsActive(user domain.User, session *sessions.Session) (bool, error)
	}

	SynchronizerImpl struct {
		store  store.Store
		logger flamingo.Logger
	}
)

func (s *SynchronizerImpl) Inject(store store.Store, logger flamingo.Logger) {
	s.store = store
	s.logger = logger
}

func (s *SynchronizerImpl) Insert(user domain.User, session *sessions.Session) error {
	concatenated := strings.Join([]string{user.Sub, time.Now().Format(time.RFC3339Nano)}, "|")
	hashBytes := sha256.Sum256([]byte(concatenated))
	hash := fmt.Sprintf("%x", hashBytes)
	session.Values[hashKey] = hash

	if err := s.store.DestroySessionsForUser(user); err != nil {
		s.logger.WithField("destroySession", "failed").Error(err.Error())
	}

	return s.store.SetHashAndSessionIdForUser(user, hash, session.ID)
}

func (s *SynchronizerImpl) IsActive(user domain.User, session *sessions.Session) (bool, error) {
	hash, err := s.store.GetHashByUser(user)
	if err != nil {
		return false, err
	}

	return hash == session.Values[hashKey], nil
}
