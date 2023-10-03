package svc

// ErrValidation is returned when input is invalid
type ErrValidation struct {
	msg string
}

func (e *ErrValidation) Error() string {
	return e.msg
}

func NewErrValidation(msg string) error {
	return &ErrValidation{msg: msg}
}

// ErrServerError is returned when server encounters an error
type ErrServerError struct {
	msg   string
	cause error
}

func (e *ErrServerError) Unwrap() error {
	return e.cause
}

func (e *ErrServerError) Error() string {
	return e.msg
}

func NewErrServerError(msg string, cause error) error {
	return &ErrServerError{msg: msg, cause: cause}
}

// ErrConflict is returned when there is a conflict
type ErrConflict struct {
	msg string
}

func (e *ErrConflict) Error() string {
	return e.msg
}

func NewErrConflict(msg string) error {
	return &ErrConflict{msg: msg}
}

// ErrNotFound is returned when resource is not found
type ErrNotFound struct {
	msg string
}

func (e *ErrNotFound) Error() string {
	return e.msg
}

func NewErrNotFound(msg string) error {
	return &ErrNotFound{msg: msg}
}
