package eu

type EasyUseError struct {
	Msg string
}

func (receiver EasyUseError) Error() string {
	return receiver.Msg
}

func NewEasyUseError(Msg string) *EasyUseError {
	return &EasyUseError{Msg: Msg}
}
