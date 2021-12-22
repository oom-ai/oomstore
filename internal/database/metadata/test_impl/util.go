package test_impl

import (
	"context"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata"
)

type PrepareStoreFn func(t *testing.T) (context.Context, metadata.Store)

type DestroyStoreFn func()
