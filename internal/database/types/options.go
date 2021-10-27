package options

import "github.com/onestore-ai/onestore/pkg/onestore/types"

type MultiGetOnlineFeatureValuesOpt struct {
	DataTable  string
	EntityName string
	RevisionId int32
	EntityKeys []string
	Features   []*types.Feature
}

type CreateFeatureOpt struct {
	types.CreateFeatureOpt
	ValueType string
}

type GetFeatureValuesStreamOpt struct {
	DataTable    string
	EntityName   string
	FeatureNames []string
	Limit        *uint64
}

type InsertRevisionOpt struct {
	Revision    int64
	GroupName   string
	DataTable   string
	Description string
}
