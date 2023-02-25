package security

type SecurityError struct {
	Msg string
}

func (receiver SecurityError) Error() string {
	return receiver.Msg
}

func NewSecurityError(Msg string) *SecurityError {
	return &SecurityError{Msg: Msg}
}
