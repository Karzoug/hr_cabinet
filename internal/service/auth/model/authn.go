package model

// AuthnDAO - authn data for database exchange.
type AuthnDAO struct {
	UserID       int    `db:"user_id"`
	RoleID       int    `db:"role_id"`
	PasswordHash string `db:"password_hash"`
}
