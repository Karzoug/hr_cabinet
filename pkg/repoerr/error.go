package repoerr

import "errors"

var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrRecordAlreadyExist = errors.New("record already exist")
)
