package offline

import "github.com/oom-ai/oomstore/pkg/oomstore/types"

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
	Features       types.RichFeatureList
}

type ImportOpt struct {
	types.ImportBatchFeaturesOpt
	Entity   *types.Entity
	Features types.FeatureList
	Header   []string
}
