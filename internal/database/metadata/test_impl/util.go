package test_impl

import (
	"context"

	"github.com/ethhte88/oomstore/internal/database/metadata"
)

type PrepareStoreFn func() (context.Context, metadata.Store)
