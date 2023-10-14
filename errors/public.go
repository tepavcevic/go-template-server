package errors

// Public wraps the original error with a new error that has a
// `Public() string` method that returns a message acceptable
// to the end user. This error can be unwrapped using the `errors`
// package approach with `Unwrap()`
func Public(err error, msg string) error {
	return publicError{err, msg}
}

type publicError struct {
	err error
	msg string
}

func (pe publicError) Error() string {
	return pe.err.Error()
}

func (pe publicError) Public() string {
	return pe.msg
}

func (pe publicError) Unwrap() error {
	return pe.err
}
