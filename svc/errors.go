package svc

import "errors"

var (
	ErrServerError = errors.New("server error")
	ErrConflict    = errors.New("conflict")
)

type ErrValidationFailed struct {
	msg string
}

func (e *ErrValidationFailed) Error() string {
	return e.msg
}

func (e *ErrValidationFailed) Is(err error) bool {
	_, ok := err.(*ErrValidationFailed)
	return ok
}

func NewErrValidationFailed(msg string) *ErrValidationFailed {
	return &ErrValidationFailed{msg: msg}
}
