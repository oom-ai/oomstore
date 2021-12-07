package offline

import (
	"encoding/csv"

	"github.com/ethhte88/oomstore/internal/database/metadata"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

type ExportOpt struct {
	DataTable  string
	EntityName string
	Features   types.FeatureList
	Limit      *uint64
}

type JoinOpt struct {
	Entity           types.Entity
	EntityRows       <-chan types.EntityRow
	FeatureMap       map[string]types.FeatureList
	RevisionRangeMap map[string][]*metadata.RevisionRange
	ValueNames       []string
}

type JoinOneGroupOpt struct {
	GroupName           string
	Features            types.FeatureList
	RevisionRanges      []*metadata.RevisionRange
	Entity              types.Entity
	EntityRowsTableName string
	ValueNames          []string
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
