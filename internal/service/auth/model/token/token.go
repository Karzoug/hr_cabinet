package token

import (
	"errors"
	"time"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type Data struct {
	UserID int
	RoleID int
}

// Payload содержит полезную нагрузку для токена.
type Payload struct {
	Data      Data      `json:"data"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload создаёт объект Payload.
func NewPayload(data Data, duration time.Duration) (*Payload, error) {
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

func (p *Payload) GetData() Data {
	return p.Data
}
