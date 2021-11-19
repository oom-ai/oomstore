package errdefs

import "errors"

func IsNotFound(err error) bool {
	return errors.As(err, &errNotFound{})
}
