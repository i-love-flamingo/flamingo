package application

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"strconv"

	"flamingo.me/flamingo/core/captcha/domain"
	"github.com/dchest/captcha"
	"github.com/pkg/errors"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

type (
	// Generator to create captchas
	Generator struct {
		encryptionKey [32]byte
	}
)

// Inject dependencies
func (g *Generator) Inject(
	config *struct {
		EncryptionPassPhrase string `inject:"config:captcha.encryptionPassPhrase"`
	},
) {
	// values are standard recommendation of scrypt docu
	salt := captcha.RandomDigits(8)
	key, err := scrypt.Key([]byte(config.EncryptionPassPhrase), salt, 32768, 8, 1, 32)
	if err != nil {
		panic(err)
	}
	copy(g.encryptionKey[:], key[:32])
}

// NewCaptcha generates a new digit-only captcha with the given length
func (g *Generator) NewCaptcha(length int) *domain.Captcha {
	digits := captcha.RandomDigits(length)

	var solution string

	for _, d := range digits {
		solution += strconv.Itoa(int(d))
	}

	return g.NewCaptchaBySolution(solution)
}

// NewCaptchaByHash recreates a captcha by the given hash.
// The hash must be generated with the same encryption key, normally by the same instance of Generator
func (g *Generator) NewCaptchaByHash(hash string) (*domain.Captcha, error) {
	encrypted, err := base64.URLEncoding.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	var decryptNonce [24]byte
	copy(decryptNonce[:], encrypted[:24])
	decrypted, ok := secretbox.Open(nil, encrypted[24:], &decryptNonce, &g.encryptionKey)
	if !ok {
		return nil, errors.New("invalid key")
	}

	return &domain.Captcha{
		Solution: string(decrypted),
		Hash:     hash,
	}, nil
}

// NewCaptchaBySolution creates a new captcha containing the given solution.
// However, multiple calls with the same solution will have different hashes because the encryption nonce is selected by random
func (g *Generator) NewCaptchaBySolution(solution string) *domain.Captcha {
	solutionBytes := []byte(solution)

	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}

	// nonce is stored in first 24 bytes of the encrypted string
	encrypted := secretbox.Seal(nonce[:], solutionBytes, &nonce, &g.encryptionKey)

	hash := base64.URLEncoding.EncodeToString(encrypted)

	return &domain.Captcha{
		Solution: solution,
		Hash:     hash,
	}
}
