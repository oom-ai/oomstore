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

type ExportRecord []interface{}

func (r ExportRecord) EntityKey() string {
	return r[0].(string)
}

func (r ExportRecord) ValueAt(i int) interface{} {
	return r[i+1]
}

type EntityRow struct {
	EntityKey string
	UnixMilli int64
	Values    []string
}

type JoinResult struct {
	Header []string
	Data   <-chan []interface{}
}

type ExportResult struct {
	Header []string
	Data   <-chan ExportRecord
	error  <-chan error
}

func NewExportResult(header []string, data <-chan ExportRecord, error <-chan error) *ExportResult {
	return &ExportResult{
		Header: header,
		Data:   data,
		error:  error,
	}
}

// ATTENTION: call this method only after you consume all elements
// from Data channel; otherwise, it will block the Data channel.
func (e *ExportResult) CheckStreamError() error {
	if e == nil {
		return nil
	}
	if e.error != nil {
		return <-e.error
	}
	return nil
}

func (e *ExportResult) GetErrorChannel() <-chan error {
	return e.error
}

type DataTableSchema struct {
	Fields []DataTableFieldSchema
}

type DataTableFieldSchema struct {
	Name      string
	ValueType ValueType
}
