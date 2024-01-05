package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/Employee-s-file-cabinet/backend/internal/service/auth/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

const (
	selectAuthnData = `
select users.id as user_id, role_id, password_hash
from users
join authorizations a on users.id = a.user_id
where work_email=$1;`
)

func (s *storage) Get(ctx context.Context, login string) (model.AuthnDAO, error) {
	rows, err := s.DB.Query(ctx, selectAuthnData, login)
	if err != nil {
		return model.AuthnDAO{}, err
	}
	authnData, err := pgx.CollectExactlyOneRow[model.AuthnDAO](rows, pgx.RowToStructByNameLax[model.AuthnDAO])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return authnData, repoerr.ErrRecordNotFound
		}
		return authnData, err
	}
	return authnData, nil
}
