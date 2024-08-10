package errors

import "fmt"

var (
	ErrNotFound = fmt.Errorf("resource not found")
)

type Error struct {
	Code    int
	Message string
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
