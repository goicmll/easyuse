package habits

type HabitError struct {
	Msg string
}

func (receiver HabitError) Error() string {
	return receiver.Msg
}

func NewHabitError(Msg string) *HabitError {
	return &HabitError{Msg: Msg}
}
