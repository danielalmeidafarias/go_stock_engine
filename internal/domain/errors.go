package domain

type ErrorCode int

const (
	ErrNotFound ErrorCode = iota + 1
	ErrConflict
	ErrBadRequest
	ErrInternal
)

type Error struct {
	ErrCode ErrorCode
	Message string
}

func (err *Error) Error() string {
	return err.Message
}

func NewError(message string, errCode ErrorCode) *Error {
	return &Error{
		ErrCode: errCode,
		Message: message,
	}
}
