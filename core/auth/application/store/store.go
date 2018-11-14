package store

import (
	"flamingo.me/flamingo/core/auth/domain"
)

type (
	Store interface {
		SetHashForUser(user domain.User, hash string) error
		GetHashByUser(user domain.User) (string, error)
	}
)
