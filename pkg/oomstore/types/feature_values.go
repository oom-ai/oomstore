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
