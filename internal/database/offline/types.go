package offline

import (
	"bufio"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type ExportOpt struct {
	SnapshotTable string
	EntityName    string
	Features      types.FeatureList
	Limit         *uint64
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
	Category            types.Category
	Features            types.FeatureList
	RevisionRanges      []*metadata.RevisionRange
	Entity              types.Entity
	EntityRowsTableName string
	ValueNames          []string
}

type ImportOpt struct {
	Entity            *types.Entity
	Features          types.FeatureList
	Header            []string
	Revision          *int64
	SnapshotTableName string
	Source            *CSVSource
}

type CSVSource struct {
	Reader    *bufio.Reader
	Delimiter string
}
