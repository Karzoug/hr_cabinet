package password_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/Employee-s-file-cabinet/backend/internal/service/auth/model/password"
	"github.com/Employee-s-file-cabinet/backend/pkg/rndtest"
)

func TestPassword(t *testing.T) {
	pass := password.New()
	t.Run("object pass created", func(t *testing.T) {
		require.NotNil(t, pass)
	})

	rndPassword := randomPassword()
	hashedPassword1, err := pass.Hash(rndPassword)

	t.Run("correct hashing", func(t *testing.T) {
		require.NoError(t, err)
		require.NotEmpty(t, hashedPassword1)
	})

	t.Run("correct check", func(t *testing.T) {
		err = pass.Check(rndPassword, hashedPassword1)
		require.NoError(t, err)
	})

	t.Run("wrong password check", func(t *testing.T) {
		wrongPassword := randomPassword()
		err = pass.Check(wrongPassword, hashedPassword1)
		require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
	})

	t.Run("check that hashes of right and wrong password not equal", func(t *testing.T) {
		hashedPassword2, err := pass.Hash(rndPassword)
		require.NoError(t, err)
		require.NotEmpty(t, hashedPassword1)
		require.NotEqual(t, hashedPassword1, hashedPassword2)
	})
}

// Login генерирует случайный логин.
func login() string {
	return rndtest.String(rndtest.Int(6, 12))
}

// Password генерирует случайный пароль.
func randomPassword() string {
	return rndtest.String(rndtest.Int(8, 24))
}
