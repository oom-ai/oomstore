package online

import (
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

type GetOpt struct {
	EntityName  string
	RevisionId  int32
	EntityKey   string
	FeatureList typesv2.FeatureList
}

type MultiGetOpt struct {
	EntityName  string
	RevisionId  int32
	EntityKeys  []string
	FeatureList typesv2.FeatureList
}

type ImportOpt struct {
	Features typesv2.FeatureList
	Revision *typesv2.Revision
	Entity   *typesv2.Entity
	Stream   <-chan *types.RawFeatureValueRecord
}
