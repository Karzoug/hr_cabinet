package service

type errorStatus uint8

const (
	InvalidArgument errorStatus = iota
	NotFound
	AlreadyExists
	NotModified
	Conflict
	PermissionDenied
	Unauthenticated
)

type Error struct {
	Status errorStatus
	text   string
}

func NewError(status errorStatus, text string) *Error {
	return &Error{
		Status: status,
		text:   text,
	}
}

func (err *Error) Error() string {
	return err.text
}
