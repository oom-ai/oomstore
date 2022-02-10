package types

type Category = string

const (
	CategoryBatch  Category = "batch"
	CategoryStream Category = "stream"
)

type TableType = string

const (
	TableBatchSnapshot  TableType = "batch_snapshot"
	TableStreamSnapshot TableType = "stream_snapshot"
	TableStreamCdc      TableType = "stream_cdc"
)

type ExportRecord struct {
	Record []interface{}
	Error  error
}

func (r ExportRecord) EntityKey() string {
	return r.Record[0].(string)
}

func (r ExportRecord) ValueAt(i int) interface{} {
	return r.Record[i+1]
}

type ExportResult struct {
	Header []string
	Data   <-chan ExportRecord
}

func NewExportResult(header []string, data <-chan ExportRecord) *ExportResult {
	return &ExportResult{
		Header: header,
		Data:   data,
	}
}

type EntityRow struct {
	EntityKey string
	UnixMilli int64
	Values    []string
	Error     error
}

type JoinRecord struct {
	Record []interface{}
	Error  error
}

type JoinResult struct {
	Header []string
	Data   <-chan JoinRecord
}

type DataTableTimeRange struct {
	MinUnixMilli *int64 `db:"min_unix_milli"`
	MaxUnixMilli *int64 `db:"max_unix_milli"`
}

type DataTableSchema struct {
	Fields    []DataTableFieldSchema
	TimeRange DataTableTimeRange
}

type DataTableFieldSchema struct {
	Name      string
	ValueType ValueType
}
