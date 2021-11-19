package types

const (
	BatchFeatureCategory  = "batch"
	StreamFeatureCategory = "stream"
)

type FeatureKV struct {
	FeatureName  string
	FeatureValue interface{}
}

func NewFeatureKV(name string, value interface{}) FeatureKV {
	return FeatureKV{
		FeatureName:  name,
		FeatureValue: value,
	}
}

type FeatureDataSet map[string][]FeatureKV

func NewFeatureDataSet() FeatureDataSet {
	return make(map[string][]FeatureKV)
}

type ExportRecord struct {
	Record []interface{}
	Error  error
}

func (r *ExportRecord) EntityKey() string {
	return r.Record[0].(string)
}

func (r *ExportRecord) ValueAt(i int) interface{} {
	return r.Record[i+1]
}

type EntityRow struct {
	EntityKey string `db:"entity_key"`
	UnixTime  int64  `db:"unix_time"`
}

type JoinResult struct {
	Header []string
	Data   <-chan []interface{}
}
