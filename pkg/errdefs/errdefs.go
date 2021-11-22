package errdefs

type errNotFound struct{ error }

func (e errNotFound) Unwrap() error {
	return e.error
}

func NotFound(err error) error {
	if err == nil || IsNotFound(err) {
		return err
	}
	return errNotFound{err}
}
