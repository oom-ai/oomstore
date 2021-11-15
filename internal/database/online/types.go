package online

import (
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type GetOpt struct {
	Entity      *types.Entity
	RevisionID  int32
	EntityKey   string
	FeatureList types.FeatureList
}

type MultiGetOpt struct {
	Entity      *types.Entity
	RevisionID  int32
	EntityKeys  []string
	FeatureList types.FeatureList
}

type ImportOpt struct {
	Revision    *types.Revision
	Entity      *types.Entity
	Stream      <-chan *types.RawFeatureValueRecord
	FeatureList types.FeatureList
}
