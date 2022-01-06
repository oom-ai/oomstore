package cmd

import (
	"context"
	"log"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/types/apply"
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
		log.Fatalf("failed opening OomStore: %+v", err)
	}
	return store
}

func stringPtr(s string) *string {
	return &s
}

func groupsToApplyGroupItems(ctx context.Context, store *oomstore.OomStore, groups types.GroupList) (*apply.GroupItems, error) {
	// TODO: Use group ids to filter, rather than taking them all out
	features, err := store.ListFeature(ctx, types.ListFeatureOpt{})
	if err != nil {
		return nil, err
	}
	return apply.FromGroupList(groups, features), nil
}
