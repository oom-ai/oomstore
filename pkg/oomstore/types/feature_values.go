package types

type FeatureValues struct {
	EntityName      string
	EntityKey       string
	FeatureNames    []string
	FeatureValueMap map[string]interface{}
}

func (fv *FeatureValues) FeatureValueSlice() []interface{} {
	values := make([]interface{}, 0, len(fv.FeatureNames))
	for _, name := range fv.FeatureNames {
		values = append(values, fv.FeatureValueMap[name])
	}
	return values
}
