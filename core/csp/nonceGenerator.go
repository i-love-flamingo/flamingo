package csp

import (
	"github.com/satori/go.uuid"
)

type (
	// NonceGenerator is an interface to generate a nonce
	NonceGenerator interface {
		GenerateNonce() string
	}

	UuidGenerator struct{}
)

// generateNonce generates a nonce
func (*UuidGenerator) GenerateNonce() string {
	return uuid.NewV4().String()
}
