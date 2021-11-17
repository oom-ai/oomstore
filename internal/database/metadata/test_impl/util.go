package test_impl

import (
	"context"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata"
)

type PrepareStoreRuntimeFunc func(t *testing.T) (context.Context, metadata.Store)
