package auth

import "errors"

var (
	ErrForbidden      = errors.New("forbidden")
	ErrTokenIsInvalid = errors.New("access token is missing or invalid")
	ErrUserNotAllowed = errors.New("user is not allowed to access")
)
