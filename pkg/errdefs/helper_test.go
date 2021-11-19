package errdefs_test

import (
	"errors"
	"testing"

	"github.com/oom-ai/oomstore/pkg/errdefs"
)

var errTest = errors.New("this is a test")

func TestErrNotFound(t *testing.T) {
	if errdefs.IsNotFound(errTest) {
		t.Fatalf("did not expect not found error, got %T", errTest)
	}

	e := errdefs.NotFound(errTest)

	if !errdefs.IsNotFound(e) {
		t.Fatalf("expected not found error, got: %T", e)
	}

	if !errors.Is(e, errTest) {
		t.Fatalf("expected not found error to match errTest")
	}
}
