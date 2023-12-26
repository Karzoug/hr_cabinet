package token

import (
	"errors"
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/model"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

// Payload содержит полезную нагрузку для токена.
type Payload struct {
	Data      model.TokenData `json:"data"`
	ExpiredAt time.Time       `json:"expired_at"`
}

// NewPayload создаёт объект Payload.
func NewPayload(data model.TokenData, duration time.Duration) (*Payload, error) {
	return &Payload{
		Data:      data,
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

// Valid - проверяет валидность токена.
func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func (p *Payload) GetData() model.TokenData {
	return p.Data
}
