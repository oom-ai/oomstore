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
	if !debug {
		return err
	}

	return errors.WithStack(err)
}

func Cause(err error) error {
	return errors.Cause(err)
}
