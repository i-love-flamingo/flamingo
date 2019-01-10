package application

import (
	"strconv"
)

type (
	// PseudoStore implements the Store interface of github.com/dchest/captcha by using regeneration by hash
	PseudoStore struct {
		generator *Generator
	}
)

// Inject dependencies
func (s *PseudoStore) Inject(
	g *Generator,
) {
	s.generator = g
}

// Set sets the digits for the captcha id.
func (s *PseudoStore) Set(_ string, _ []byte) {
	return
}

// Get returns the original digits from the given hash
func (s *PseudoStore) Get(id string, _ bool) []byte {
	c, err := s.generator.NewCaptchaByHash(id)
	if err != nil {
		return nil
	}

	var digits []byte

	for _, char := range c.Solution {
		d, err := strconv.Atoi(string(char))
		if err != nil {
			return nil
		}

		digits = append(digits, uint8(d))
	}

	return digits
}
