package cmd

import (
	"context"
	"log"

	"github.com/ethhte88/oomstore/pkg/oomstore"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

const (
	CSV        = "csv"
	ASCIITable = "ascii_table"
	Column     = "column"
	YAML       = "yaml"
)

func mustOpenOomStore(ctx context.Context, opt types.OomStoreConfig) *oomstore.OomStore {
	store, err := oomstore.Open(ctx, oomStoreCfg)
	if err != nil {
		log.Fatalf("failed opening OomStore: %v", err)
	}
	return store
}

func stringPtr(s string) *string {
	return &s
}
