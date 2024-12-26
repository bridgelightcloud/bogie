package errors

type Error struct {
	msg  string
	errs []error
}

func (e Error) Error() string {
	return e.msg
}

func (e Error) Unwrap() []error {
	return e.errs
}
