package offline

import (
	"bufio"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type ExportOpt struct {
	SnapshotTables map[int]string
	CdcTables      map[int]string
	Features       map[int]types.FeatureList
	UnixMilli      int64
	EntityName     string
	Limit          *uint64
}

type RevisionRange struct {
	MinRevision   int64
	MaxRevision   int64
	SnapshotTable string
	CdcTable      string
}

type JoinOpt struct {
	Entity           types.Entity
	EntityRows       <-chan types.EntityRow
	FeatureMap       map[string]types.FeatureList
	RevisionRangeMap map[string][]*RevisionRange
	ValueNames       []string
}

type JoinOneGroupOpt struct {
	GroupName           string
	Category            types.Category
	Features            types.FeatureList
	RevisionRanges      []*RevisionRange
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
	NoPK              bool // TODO: to import cdc data temporarily for testing
}

type PushOpt struct {
	GroupID      int
	Revision     int64
	EntityName   string
	FeatureNames []string
	Records      []types.StreamRecord
}

type CSVSource struct {
	Reader    *bufio.Reader
	Delimiter string
}

type SnapshotOpt struct {
	Group        *types.Group
	Features     types.FeatureList
	Revision     int64
	PrevRevision int64
}

type CreateTableOpt struct {
	TableName string
	Entity    *types.Entity
	Features  types.FeatureList
	TableType types.TableType
}
