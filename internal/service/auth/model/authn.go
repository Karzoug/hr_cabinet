package model

// AuthnDAO - authn data for database exchange.
type AuthnDAO struct {
	UserID       string `db:"user_id"`
	RoleID       string `db:"role_id"`
	PasswordHash string `db:"password_hash"`
}
