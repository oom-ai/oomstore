package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	CSV        = "csv"
	ASCIITable = "ascii_table"
	Column     = "column"
	YAML       = "yaml"
)

const (
	MetadataFieldTruncateAt = 40
)

func mustOpenOomStore(ctx context.Context, opt types.OomStoreConfig) *oomstore.OomStore {
	store, err := oomstore.Open(ctx, oomStoreCfg)
	if err != nil {
		exitf("failed opening OomStore: %+v", err)
	}
	return store
}

func stringPtr(s string) *string {
	return &s
}

func exitf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	os.Exit(1)
}

func exit(a ...interface{}) {
	msg := fmt.Sprint(a...)
	fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	os.Exit(1)
}
