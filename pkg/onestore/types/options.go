package types

type OneStoreOpt struct {
	Host      string
	Port      string
	User      string
	Pass      string
	Workspace string
}

type CreateFeatureOpt struct {
	FeatureName string
	GroupName   string
	ValueType   string
	Description string
}

type ListFeatureOpt struct {
	EntityName *string
	GroupName  *string
}

type UpdateFeatureOpt struct {
	FeatureName    string
	NewDescription string
}

type CreateFeatureGroupOpt struct {
	Name        string
	EntityName  string
	Category    string
	Description string
}

type WalkFeatureValuesFunc = func(values []interface{}) error

type WalkFeatureValuesOpt struct {
	FeatureGroup          FeatureGroup
	FeatureNames          []string
	Limit                 *uint64
	WalkFeatureValuesFunc WalkFeatureValuesFunc
}

type CreateBatchFeatureOpt struct {
	FeatureName string
	GroupName   string
	ValueType   string
	Description string
}
