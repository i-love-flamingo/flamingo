package store

import (
	"flamingo.me/flamingo/core/auth/domain"
)

type (
	Nil struct{}
)

func (s *Nil) SetHashForUser(user domain.User, hash string) error {
	return nil
}

func (s *Nil) GetHashByUser(user domain.User) (string, error) {
	return "", nil
}
