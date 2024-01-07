package token_test

import (
	"strings"
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
		maker, err = token.NewPasetoMaker(gofakeit.Lexify(strings.Repeat("?", 32)), duration)
		require.NoError(t, err)
	})

	var testToken string
	t.Run("create token", func(t *testing.T) {
		testToken, err = maker.Create(data)
		require.NoError(t, err)
		require.NotEmpty(t, testToken)
	})

	t.Run("verify token", func(t *testing.T) {
		payload, err := maker.Verify(testToken)
		require.NoError(t, err)
		require.NotEmpty(t, payload)

		require.Equal(t, data, payload.Data)
	})
}

func TestExpiredPasetoToken(t *testing.T) {
	var err error

	maker, err := token.NewPasetoMaker(gofakeit.Lexify(strings.Repeat("?", 32)), -time.Minute)
	require.NoError(t, err)

	data := token.Data{
		UserID: gofakeit.Numerify("###"),
		RoleID: gofakeit.Numerify("###"),
	}

	var testToken string
	t.Run("create expired token", func(t *testing.T) {
		testToken, err = maker.Create(data)
		require.NoError(t, err)
		require.NotEmpty(t, testToken)
	})

	t.Run("verify expired token", func(t *testing.T) {
		payload, err := maker.Verify(testToken)
		require.Error(t, err)
		require.EqualError(t, err, token.ErrExpiredToken.Error())
		require.Nil(t, payload)
	})
}
