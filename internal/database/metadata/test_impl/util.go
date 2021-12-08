package test_impl

import (
	"context"

	"github.com/ethhte88/oomstore/internal/database/metadata"
)

type PrepareStoreRuntimeFunc func() (context.Context, metadata.Store)
