package domain

type ErrorCode int

const (
	ErrNotFound ErrorCode = iota + 1
	ErrConflict
	ErrBadRequest
	ErrInternal
)

type Error struct {
	Err     error
	ErrCode ErrorCode
	Message string
}

func NewError(err error, errCode ErrorCode, message string) *Error {
	return &Error{
		Err:     err,
		ErrCode: errCode,
		Message: message,
	}
}
