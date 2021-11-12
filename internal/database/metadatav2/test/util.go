package test

import (
	"context"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
)

type PrepareStoreRuntimeFunc func(t *testing.T) (context.Context, metadatav2.Store)
