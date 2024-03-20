package usecase

type Error struct {
	msg string
}

func (e Error) Error() string {
	return e.msg
}

// NewError creates a new usecase error.
func NewError(msg string) Error {
	return Error{
		msg: msg,
	}
}
