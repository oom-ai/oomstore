package errdefs

import "errors"

func IsNotFound(err error) bool {
	return errors.As(err, &errNotFound{})
}

func IsInvalidAttribute(err error) bool {
	return errors.As(err, &errInvalidAttribute{})
}
