package options

import "github.com/onestore-ai/onestore/pkg/onestore/types"

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
