package metadata

import "github.com/onestore-ai/onestore/pkg/onestore/types"

type CreateFeatureOpt struct {
	types.CreateFeatureOpt
	ValueType string
}

type InsertRevisionOpt struct {
	Revision    int64
	GroupName   string
	DataTable   string
	Description string
}
