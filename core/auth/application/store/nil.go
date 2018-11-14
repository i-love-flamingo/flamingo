package store

import (
	"flamingo.me/flamingo/core/auth/domain"
)

type (
	Nil struct{}
)

func (s *Nil) SetHashAndSessionIdForUser(user domain.User, hash string, id string) error {
	return nil
}

func (s *Nil) GetHashByUser(user domain.User) (string, error) {
	return "", nil
}

func (s *Nil) DestroySessionsForUser(user domain.User) error {
	return nil
}
