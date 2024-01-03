package user

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrPhotoFileNotFound     = errors.New("photo file not found")
	ErrPhotoFileNotModified  = errors.New("photo file not modified")
	ErrPhotoFileSizeTooLarge = errors.New("photo file size too large")
)
