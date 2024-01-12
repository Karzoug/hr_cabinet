package postgres

import "github.com/Employee-s-file-cabinet/backend/internal/service/recovery/model"

type user struct {
	ID         int    `db:"id"`
	LastName   string `db:"lastname"`
	FirstName  string `db:"firstname"`
	MiddleName string `db:"middlename"`
	WorkEmail  string `db:"work_email"`
}

func convertUserToModelUser(user *user) model.User {
	return model.User{
		ID:         user.ID,
		LastName:   user.LastName,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		Email:      user.WorkEmail,
	}
}
