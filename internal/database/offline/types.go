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
	EntityName       string
	EntityRows       <-chan types.EntityRow
	GroupNames       []string
	FeatureMap       map[string]types.FeatureList
	RevisionRangeMap map[string][]*RevisionRange
	ValueNames       []string
}

type ImportOpt struct {
	EntityName        string
	SnapshotTableName string
	Header            []string
	Features          types.FeatureList
	Revision          *int64
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
	Group        types.Group
	Features     types.FeatureList
	Revision     int64
	PrevRevision int64
}

type CreateTableOpt struct {
	TableName  string
	EntityName string
	Features   types.FeatureList
	TableType  types.TableType
}

type TableSchemaOpt struct {
	TableName      string
	CheckTimeRange bool
}
