package application

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"time"

	"crypto/cipher"

	"net/http"

	"crypto/sha256"

	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
)

const (
	TokenName = "csrftoken"
)

type (
	Service interface {
		Generate(session *web.Session) string
		IsValid(request *web.Request) bool
	}

	ServiceImpl struct {
		secret []byte
		ttl    int

		logger flamingo.Logger
	}

	csrfToken struct {
		ID   string    `json:"id"`
		Date time.Time `json:"date"`
	}
)

func (s *ServiceImpl) Inject(l flamingo.Logger, cfg *struct {
	Secret string  `inject:"config:csrf.secret"`
	Ttl    float64 `inject:"config:csrf.ttl"`
}) {
	hash := sha256.Sum256([]byte(cfg.Secret))
	s.secret = hash[:]
	s.ttl = int(cfg.Ttl)
	s.logger = l
}

func (s *ServiceImpl) Generate(session *web.Session) string {
	token := csrfToken{
		ID:   session.ID(),
		Date: time.Now(),
	}

	body, err := json.Marshal(token)
	if err != nil {
		s.logger.WithField("csrf", "jsonMarshal").Error(err.Error())
		return ""
	}

	gcm, err := s.getGcm()
	if err != nil {
		s.logger.WithField("csrf", "newGCM").Error(err.Error())
		return ""
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		s.logger.WithField("csrf", "nonceGenerate").Error(err.Error())
		return ""
	}

	cipherText := gcm.Seal(nil, nonce, body, nil)
	cipherText = append(nonce, cipherText...)
	return hex.EncodeToString(cipherText)
}

func (s *ServiceImpl) IsValid(request *web.Request) bool {
	if request.Request().Method != http.MethodPost {
		return true
	}

	formToken, ok := request.Form1(TokenName)
	if !ok {
		return false
	}

	data, err := hex.DecodeString(formToken)
	if err != nil {
		return false
	}

	gcm, err := s.getGcm()
	if err != nil {
		return false
	}

	nonceSize := gcm.NonceSize()
	if len(data) <= nonceSize {
		return false
	}

	nonce := data[:nonceSize]
	cipherText := data[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return false
	}

	var token csrfToken
	err = json.Unmarshal(plainText, &token)
	if err != nil {
		return false
	}

	if request.Session().ID() != token.ID {
		return false
	}

	if time.Now().Add(time.Duration(-s.ttl) * time.Second).After(token.Date) {
		return false
	}

	return true
}

func (s *ServiceImpl) getGcm() (cipher.AEAD, error) {
	block, err := aes.NewCipher(s.secret)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm, nil
}
