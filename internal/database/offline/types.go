package offline

import "github.com/onestore-ai/onestore/pkg/onestore/types"

type ExportOpt struct {
	DataTable    string
	EntityName   string
	FeatureNames []string
	Limit        *uint64
}

type JoinOpt struct {
	Entity         *types.Entity
	EntityRows     []types.EntityRow
	RevisionRanges []*types.RevisionRange
	Features       []*types.RichFeature
}

type ImportOpt struct {
	types.ImportBatchFeaturesOpt
	Entity   *types.Entity
	Features []*types.Feature
	Header   []string
}
