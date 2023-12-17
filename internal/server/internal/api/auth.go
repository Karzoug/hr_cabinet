package api

type AccessPermission uint64

const (
	UserCreateAccess AccessPermission = 1 << iota
	UserReadAccess
	UserUpdateAccess
	UserDeleteAccess
)
