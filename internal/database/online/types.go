package online

import (
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type GetOpt struct {
	EntityKey  string
	Group      types.Group
	Features   types.FeatureList
	RevisionID *int
}

type MultiGetOpt struct {
	EntityKeys []string
	Group      types.Group
	Features   types.FeatureList
	RevisionID *int
}

type ImportOpt struct {
	Group        types.Group
	Features     types.FeatureList
	ExportStream <-chan types.ExportRecord
	ExportError  <-chan error
	RevisionID   *int
}

type PushOpt struct {
	Entity        *types.Entity
	EntityKey     string
	GroupID       int
	Features      types.FeatureList
	FeatureValues []interface{}
}

type PrepareStreamTableOpt struct {
	Entity *types.Entity

	GroupID int

	// Feature is not nil to add a new row to the stream table;
	// otherwise it means the stream table will be created.
	Feature *types.Feature
}
