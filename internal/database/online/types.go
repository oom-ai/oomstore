package online

import (
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type GetOpt struct {
	Entity     *types.Entity
	EntityKey  string
	RevisionID *int
	Group      *types.Group
	Features   types.FeatureList
}

type MultiGetOpt struct {
	Entity     *types.Entity
	EntityKeys []string
	RevisionID *int
	Group      *types.Group
	Features   types.FeatureList
}

type ImportOpt struct {
	Group        types.Group
	Features     types.FeatureList
	Revision     *types.Revision
	ExportStream <-chan types.ExportRecord
	ExportError  <-chan error
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
