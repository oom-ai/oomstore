package online

import (
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type GetOpt struct {
	Entity      *types.Entity
	RevisionID  int
	EntityKey   string
	FeatureList types.FeatureList
}

type MultiGetOpt struct {
	Entity      *types.Entity
	RevisionID  int
	EntityKeys  []string
	FeatureList types.FeatureList
}

type ImportOpt struct {
	Revision     *types.Revision
	Entity       *types.Entity
	ExportStream <-chan types.ExportRecord
	ExportError  <-chan error
	FeatureList  types.FeatureList
}
