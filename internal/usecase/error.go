package usecase

type Error struct {
	// from      string
	// tagNeeded bool
	msg string
}

func (e Error) Error() string {
	// if e.tagNeeded {
	// 	return fmt.Sprintf("@%s, %s", e.from, e.msg)
	// } else {
	return e.msg
	// }
}

// NewError creates a new service error.
// If from is not empty, the error will mention the from user.
func NewError(msg string) Error {
	err := Error{
		msg: msg,
		//from: from,
	}
	// if from != "" {
	// 	err.tagNeeded = true
	// }
	return err
}
