package easyuse

type EasyUseError struct {
	Msg string
}

func (e EasyUseError) Error() string {
	return e.Msg
}

func NewEasyUseError(Msg string) *EasyUseError {
	return &EasyUseError{Msg: Msg}
}
