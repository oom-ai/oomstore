package online

import (
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type GetOpt struct {
	Entity     *types.Entity
	RevisionID int
	EntityKey  string
	Features   types.FeatureList
}

type MultiGetOpt struct {
	Entity     *types.Entity
	RevisionID int
	EntityKeys []string
	Features   types.FeatureList
}

type ImportOpt struct {
	Revision     *types.Revision
	Entity       *types.Entity
	ExportStream <-chan types.ExportRecord
	ExportError  <-chan error
	Features     types.FeatureList
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
