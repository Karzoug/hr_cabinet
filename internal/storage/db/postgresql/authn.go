package postgresql

import (
	"context"
)

// merged --->
// TODO: Перенести в repo
const (

	selectUserID = `
select users.id as user_id
from users
join authorizations a on users.id = a.user_id
where work_email=$1;`

	updatePass = `
update authorizations
set password_hash=$1
where user_id=$2;`
)

func (s *storage) ExistEmployee(ctx context.Context, workEmail string) (bool, int, error) {
	row := s.QueryRowContext(ctx, selectUserID, workEmail)
	var userID int
	err := row.Scan(&userID)
	if err != nil || userID == 0 {
		return false, 0, err
	}
	return true, userID, nil
}

func (s *storage) ChangePass(ctx context.Context, userID int, hash string) error {
	_, err := s.ExecContext(ctx, updatePass, hash)
	return err
}
// <---