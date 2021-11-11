package offline

import (
	"encoding/csv"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

type ExportOpt struct {
	DataTable    string
	EntityName   string
	FeatureNames []string
	Limit        *uint64
}

type JoinOpt struct {
	Entity           typesv2.Entity
	EntityRows       <-chan types.EntityRow
	FeatureMap       map[string]typesv2.FeatureList
	RevisionRangeMap map[string][]*metadatav2.RevisionRange
}

type JoinOneFeatureGroupOpt struct {
	GroupName           string
	Features            typesv2.FeatureList
	RevisionRanges      []*metadatav2.RevisionRange
	Entity              typesv2.Entity
	EntityRowsTableName string
}

type ImportOpt struct {
	GroupName string
	Entity    *typesv2.Entity
	Features  typesv2.FeatureList
	Header    []string
	Revision  *int64

	// CsvReader must not contain header
	CsvReader *csv.Reader
}
