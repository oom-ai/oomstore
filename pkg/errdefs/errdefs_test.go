package errdefs_test

import (
	"errors"
	"testing"

	"github.com/oom-ai/oomstore/pkg/errdefs"
)

var errTest = errors.New("this is a test")

func TestErrNotFound(t *testing.T) {
	if err := errdefs.NotFound(nil); err != nil {
		t.Fatalf("want nil, but got %v", err)
	}

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

func TestErrInvalidAttribute(t *testing.T) {
	if err := errdefs.InvalidAttribute(nil); err != nil {
		t.Fatalf("want nil, but got %v", err)
	}

	if errdefs.IsInvalidAttribute(errTest) {
		t.Fatalf("did not expect not found error, got %T", errTest)
	}

	e := errdefs.NotFound(errTest)

	if errdefs.IsInvalidAttribute(e) {
		t.Fatalf("did not expect invalid parameter error, got %T", e)
	}

	e = errdefs.InvalidAttribute(errTest)
	if !errdefs.IsInvalidAttribute(e) {
		t.Fatalf("expected invalid parameter error, got: %T", e)
	}

	if !errors.Is(e, errTest) {
		t.Fatalf("expected invalid parameter error to match errTest")
	}
}
