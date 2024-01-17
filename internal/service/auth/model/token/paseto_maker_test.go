package token_test

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"

	"github.com/Employee-s-file-cabinet/backend/internal/service/auth/model/token"
)

func TestPasetoMaker(t *testing.T) {
	data := token.Data{
		UserID: gofakeit.Numerify("###"),
		RoleID: gofakeit.Numerify("###"),
	}
	duration := time.Minute

	var maker *token.PasetoMaker
	var err error
	t.Run("create token maker", func(t *testing.T) {
		_, private, err := ed25519.GenerateKey(nil)
		require.NoError(t, err)
		maker, err = token.NewPasetoMaker(hex.EncodeToString(private), duration)
		require.NoError(t, err)
	})

	var testToken, testSign string
	t.Run("create token", func(t *testing.T) {
		testToken, testSign, err = maker.Create(data)
		require.NoError(t, err)
		require.NotEmpty(t, testToken)
		require.NotEmpty(t, testSign)
	})

	t.Run("verify token", func(t *testing.T) {
		payload, err := maker.Verify(testToken, testSign)
		require.NoError(t, err)
		require.NotEmpty(t, payload)

		require.Equal(t, data, payload.Data)
	})
}

func TestExpiredPasetoToken(t *testing.T) {
	var err error

	_, private, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)
	maker, err := token.NewPasetoMaker(hex.EncodeToString(private), -time.Minute)
	require.NoError(t, err)

	data := token.Data{
		UserID: gofakeit.Numerify("###"),
		RoleID: gofakeit.Numerify("###"),
	}

	var testToken, testSign string
	t.Run("create expired token", func(t *testing.T) {
		testToken, testSign, err = maker.Create(data)
		require.NoError(t, err)
		require.NotEmpty(t, testToken)
	})

	t.Run("verify expired token", func(t *testing.T) {
		payload, err := maker.Verify(testToken, testSign)
		require.Error(t, err)
		require.EqualError(t, err, token.ErrExpiredToken.Error())
		require.Nil(t, payload)
	})
}
