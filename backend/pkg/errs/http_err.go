package errs

import "errors"

var (
	ErrInvalidJSON       = errors.New("invalid json")
	ErrInvalidIdentifier = errors.New("invalid identifier format")
	ErrInvalidSlug       = errors.New("invalid slug")
)

type OutErr struct {
	Code    int
	Message string
	Reason  error
}

func NewOutError(code int, msg string, reason error) *OutErr {
	return &OutErr{
		Code:    code,
		Message: msg,
		Reason:  reason,
	}
}
