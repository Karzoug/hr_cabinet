package recovery

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	serr "github.com/Employee-s-file-cabinet/backend/internal/service"
	"github.com/Employee-s-file-cabinet/backend/internal/service/recovery/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"

func (s *service) InitChangePassword(ctx context.Context, login string) error {
	user, err := s.getUser(ctx, login)
	if err != nil {
		return err
	}

	key, err := s.generateKey(user.ID)
	if err != nil {
		return err
	}

	recoveryData := model.MessageData{
		User: user,
		Key:  key,
	}

	if err = s.sendRecoveryMessage(ctx, recoveryData); err != nil {
		return err
	}

	return nil
}

func (s *service) ChangePassword(ctx context.Context, key, newPassword string) error {
	const op = "recovery service: change password"

	userID, exists := s.keyRepository.Get(key)
	if !exists {
		serr.NewError(serr.InvalidArgument, "invalid key or login")
	}

	//TODO: проверка пароля на сложность

	passHash, err := s.passwordVerificator.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.recoveryRepository.ChangePassword(ctx, userID, passHash)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// TODO: присылать сообщение об успешной смене пароля

	return nil
}

func (s *service) Check(ctx context.Context, key string) error {
	if _, exists := s.keyRepository.Get(key); !exists {
		return serr.NewError(serr.InvalidArgument, "invalid key")

	}
	return nil
}

func (s *service) getUser(ctx context.Context, login string) (*model.User, error) {
	const op = "recovery service: get user"

	user, err := s.recoveryRepository.CheckAndReturnUser(ctx, login)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.InvalidArgument, "invalid login")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *service) generateKey(userID int) (string, error) {
	const op = "recovery service: generate key"

	key, err := generateRandomString(36)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// TODO: передавать время из конфигурации
	err = s.keyRepository.Set(key, userID, s.Config.KeyLifetime)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return key, nil
}

func (s *service) sendRecoveryMessage(ctx context.Context, data model.MessageData) error {
	const op = "recovery service: send recovery message"

	// TODO: сообщение можно формировать с помощью text/template (или html/template)
	if err := s.notificationDeliverer.SendMessage(data.User.Email,
		"Завершите запрос на сброс пароля",
		fmt.Sprintf(
			`%s %s,
Для восстановления доступа к личному кабинету перейдите по ссылке:
%s/access-restore/password-reset?key=%s`,
			data.User.FirstName, data.User.LastName, s.Config.Domain, data.Key)); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func generateRandomString(n int) (string, error) {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
