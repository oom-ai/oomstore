package types

type FeatureValues struct {
	EntityName       string
	EntityKey        string
	FeatureFullNames []string
	FeatureValueMap  map[string]interface{}
}

func (fv *FeatureValues) FeatureValueSlice() []interface{} {
	values := make([]interface{}, 0, len(fv.FeatureFullNames))
	for _, name := range fv.FeatureFullNames {
		values = append(values, fv.FeatureValueMap[name])
	}
	return values
}

type StreamRecord struct {
	GroupID   int
	EntityKey string
	UnixMilli int64
	Values    []interface{}
}

func (r *StreamRecord) ToRow() []interface{} {
	row := make([]interface{}, 0, len(r.Values)+2)
	row = append(row, r.EntityKey, r.UnixMilli)
	return append(row, r.Values...)
}
