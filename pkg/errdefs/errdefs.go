package errdefs

import (
	"fmt"
)

type errNotFound struct{ error }

func (e errNotFound) Unwrap() error {
	return e.error
}

func (e errNotFound) Error() string {
	return fmt.Sprintf("%+v", e.error)
}

func NotFound(err error) error {
	if err == nil || IsNotFound(err) {
		return err
	}
	return errNotFound{err}
}

type errInvalidAttribute struct{ error }

func (e errInvalidAttribute) Unwrap() error {
	return e.error
}

func (e errInvalidAttribute) Error() string {
	return fmt.Sprintf("%+v", e.error)
}

func InvalidAttribute(err error) error {
	if err == nil || IsInvalidAttribute(err) {
		return err
	}
	return errInvalidAttribute{err}
}
