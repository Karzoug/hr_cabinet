package token

import (
	"crypto"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
)

const v2SignSize = ed25519.SignatureSize

// PasetoMaker реализация создателя токенов типа PaseTo.
type PasetoMaker struct {
	paseto     *paseto.V2
	privateKey crypto.PrivateKey
	publicKey  crypto.PublicKey
	duration   time.Duration
}

// NewPasetoMaker возвращает PasetoMaker для управления токенами.

func NewPasetoMaker(privateHexKey string, duration time.Duration) (*PasetoMaker, error) {
	const op = "create paseto maker"

	b, err := hex.DecodeString(privateHexKey)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(b) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("%s: invalid key size: must be %d characters", op, ed25519.PrivateKeySize)
	}
	privateKey := ed25519.PrivateKey(b)

	return &PasetoMaker{
		paseto:     paseto.NewV2(),
		privateKey: privateKey,
		publicKey:  privateKey.Public(),
		duration:   duration,
	}, nil
}

// Create создаёт токен для переданных данных и продолжительности.
func (m *PasetoMaker) Create(data Data) (string, string, error) {
	payload, err := NewPayload(data, m.duration)
	if err != nil {
		return "", "", err
	}

	signed, err := m.paseto.Sign(m.privateKey, payload, nil)
	if err != nil {
		return "", "", err
	}
	token := signed[:len(signed)-v2SignSize]
	sign := signed[len(signed)-v2SignSize:]

	return token, sign, nil
}

// Verify проверяет, является ли токен действительным.
func (m *PasetoMaker) Verify(token, sign string) (*Payload, error) {
	payload := &Payload{}
	err := m.paseto.Verify(token+sign, m.publicKey, payload, nil)
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
