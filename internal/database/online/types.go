package online

import "github.com/oom-ai/oomstore/pkg/onestore/types"

type GetOpt struct {
	DataTable  string
	EntityName string
	RevisionId int32
	EntityKey  string
	Features   []*types.Feature
}

type MultiGetOpt struct {
	DataTable  string
	EntityName string
	RevisionId int32
	EntityKeys []string
	Features   []*types.Feature
}

type ImportOpt struct {
	Features []*types.Feature
	Revision *types.Revision
	Entity   *types.Entity
	Stream   <-chan *types.RawFeatureValueRecord
}
