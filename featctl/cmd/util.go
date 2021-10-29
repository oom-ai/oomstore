package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func mustOpenOomStore(ctx context.Context, opt types.OomStoreOptV2) *oomstore.OomStore {
	store, err := oomstore.Open(ctx, oomStoreOpt)
	if err != nil {
		log.Fatalf("failed opening OomStore: %v", err)
	}
	return store
}
