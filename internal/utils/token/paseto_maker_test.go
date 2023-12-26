package token_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Employee-s-file-cabinet/backend/internal/model"
	"github.com/Employee-s-file-cabinet/backend/internal/utils/test"
	"github.com/Employee-s-file-cabinet/backend/internal/utils/token"
)

func TestPasetoMaker(t *testing.T) {
	data := model.TokenData{
		UserID: test.RandomInt(1, 100),
		RoleID: test.RandomInt(1, 100),
	}
	duration := time.Minute

	var maker *token.PasetoMaker
	var err error
	t.Run("create token maker", func(t *testing.T) {
		maker, err = token.NewPasetoMaker(test.RandomString(32), duration)
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

		require.Equal(t, data, payload.GetData())
	})
}

func TestExpiredPasetoToken(t *testing.T) {
	var err error

	maker, err := token.NewPasetoMaker(test.RandomString(32), -time.Minute)
	require.NoError(t, err)

	data := model.TokenData{
		UserID: test.RandomInt(1, 100),
		RoleID: test.RandomInt(1, 100),
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
