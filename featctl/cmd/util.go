package cmd

import (
	"context"
	"log"
	"strconv"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	CSV        = "csv"
	ASCIITable = "ascii_table"
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

func serializeInt16(i int16) string {
	return strconv.FormatInt(int64(i), 10)
}

func serializeInt32(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}
