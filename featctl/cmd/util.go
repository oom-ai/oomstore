package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/pkg/onestore"
	"github.com/oom-ai/oomstore/pkg/onestore/types"
)

func mustOpenOneStore(ctx context.Context, opt types.OneStoreOpt) *onestore.OneStore {
	store, err := onestore.Open(ctx, oneStoreOpt)
	if err != nil {
		log.Fatalf("failed opening OneStore: %v", err)
	}
	return store
}
