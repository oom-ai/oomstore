package online

import "github.com/oom-ai/oomstore/pkg/oomstore/types"

type GetOpt struct {
	DataTable   string
	EntityName  string
	RevisionId  int32
	EntityKey   string
	FeatureList types.FeatureList
}

type MultiGetOpt struct {
	DataTable   string
	EntityName  string
	RevisionId  int32
	EntityKeys  []string
	FeatureList types.FeatureList
}

type ImportOpt struct {
	Features types.FeatureList
	Revision *types.Revision
	Entity   *types.Entity
	Stream   <-chan *types.RawFeatureValueRecord
}
