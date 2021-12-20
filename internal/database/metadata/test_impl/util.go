package test_impl

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
)

type PrepareStoreFn func() (context.Context, metadata.Store)
