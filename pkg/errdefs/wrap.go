// A wrapper around pkg/errors.
// Prints the error call stack when the environment variable DEBUG is on, otherwise it does not print
package errdefs

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

var debug bool

func init() {
	v := os.Getenv("DEBUG")
	if v == "enabled" || v == "1" || v == "true" {
		debug = true
	}
}

func Errorf(format string, args ...interface{}) error {
	if !debug {
		return fmt.Errorf(format, args...)
	}
	return errors.Errorf(format, args...)
}

func WithStack(err error) error {
	if !debug || hasStack(err) {
		return err
	}

	return errors.WithStack(err)
}

func Cause(err error) error {
	return errors.Cause(err)
}

// not call callers() on every wrap, see https://github.com/pkg/errors/issues/75#issuecomment-574580408
func hasStack(err error) bool {
	if err == nil {
		return false
	}

	type wrapper interface {
		Unwrap() error
	}

	type tracer interface {
		StackTrace() errors.StackTrace
	}

	e := err
	for {
		if _, ok := e.(tracer); ok {
			return true
		}

		wrap, ok := e.(wrapper)
		if !ok {
			return false
		}

		e = wrap.Unwrap()
		if e == nil {
			return false
		}
	}
}
