package postgresql

import (
	"context"

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
