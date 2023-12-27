package postgresql

import (
	"context"
	"fmt"

	"github.com/Employee-s-file-cabinet/backend/internal/model"
)

const (
	selectAuthnData = `
select users.id as user_id, role_id, password_hash
from users
join authorizations a on users.id = a.user_id
where work_email=$1;`
)

func (s *storage) GetAuthnData(ctx context.Context, login string) (model.AuthnDAO, error) {
	var authnData model.AuthnDAO
	err := s.GetContext(ctx, &authnData, selectAuthnData, login)
	return authnData, err
}

func (s *storage) ExistEmployee(ctx context.Context, userID int) (bool, error) {
	//TODO: запрос на наличие пользователя
	return false, fmt.Errorf("not implemented")
}

func (s *storage) ChangePass(ctx context.Context, login, hash string) error {
	//TODO: запрос на смену пароля
	return fmt.Errorf("not implemented")
}
