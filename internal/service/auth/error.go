package auth

import serr "github.com/Employee-s-file-cabinet/backend/internal/service"

var errUnauthenticated = serr.NewError(
	serr.Unauthenticated,
	"login or password is incorrect",
)
