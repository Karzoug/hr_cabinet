package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

// PasetoMaker реализация создателя токенов типа PaseTo.
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
	duration     time.Duration
}

// NewPasetoMaker возвращает PasetoMaker для управления токенами.
func NewPasetoMaker(symmetricKey string, duration time.Duration) (*PasetoMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be %d characters", chacha20poly1305.KeySize)
	}

	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
		duration:     duration,
	}, nil
}

// Create создаёт токен для переданных данных и продолжительности.
func (m *PasetoMaker) Create(data Data) (string, error) {
	payload, err := NewPayload(data, m.duration)
	if err != nil {
		return "", err
	}

	return m.paseto.Encrypt(m.symmetricKey, payload, nil)
}

// Verify проверяет, является ли токен действительным.
func (m *PasetoMaker) Verify(in string) (*Payload, error) {
	payload := &Payload{}
	err := m.paseto.Decrypt(in, m.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

// Expires возвращает время истечения срока годности токена (начиная с текущего момента времени).
func (m *PasetoMaker) Expires() time.Time {
	return time.Now().Add(m.duration)
}
