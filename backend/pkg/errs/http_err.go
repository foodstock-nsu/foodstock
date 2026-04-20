package errs

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
