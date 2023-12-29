package postgresql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"

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
	var authnData model.AuthnDAO
	if err := pgxscan.Get(ctx, s.DB, &authnData, selectAuthnData, login); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authnData, repoerr.ErrRecordNotFound
		}
		return authnData, err
	}
	return authnData, nil
}
