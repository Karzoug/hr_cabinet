package user

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrPhotoFileSizeTooLarge = errors.New("photo file size too large")
)
