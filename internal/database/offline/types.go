package offline

import (
	"encoding/csv"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type ExportOpt struct {
	DataTable    string
	EntityName   string
	FeatureNames []string
	Limit        *uint64
}

type JoinOpt struct {
	Entity           types.Entity
	EntityRows       <-chan types.EntityRow
	FeatureMap       map[string]types.FeatureList
	RevisionRangeMap map[string][]*metadata.RevisionRange
}

type JoinOneFeatureGroupOpt struct {
	GroupName           string
	Features            types.FeatureList
	RevisionRanges      []*metadata.RevisionRange
	Entity              types.Entity
	EntityRowsTableName string
}

type ImportOpt struct {
	Entity        *types.Entity
	Features      types.FeatureList
	Header        []string
	Revision      *int64
	DataTableName string

	// CsvReader must not contain header
	CsvReader *csv.Reader
}
