package typesv2

type RawFeatureValueRecord struct {
	Record []interface{}
	Error  error
}

func (r *RawFeatureValueRecord) EntityKey() string {
	return r.Record[0].(string)
}

func (r *RawFeatureValueRecord) ValueAt(i int) interface{} {
	return r.Record[i+1]
}
