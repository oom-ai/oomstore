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

type CreateEntityOpt struct {
	Name        string
	Length      uint
	Description string
}

type CreateFeatureGroupOpt struct {
	Name        string
	EntityName  string
	Category    string
	Description string
}

type WalkFeatureValuesFunc = func(key string, values []interface{}) error

type WalkFeatureValuesOpt struct {
	FeatureGroup          FeatureGroup
	FeatureNames          []string
	Limit                 *uint64
	WalkFeatureValuesFunc WalkFeatureValuesFunc
}

type ImportBatchFeaturesOpt struct {
	GroupName   string
	Description string
	DataSource  LocalFileDataSourceOpt
}

type LocalFileDataSourceOpt struct {
	FilePath  string
	Separator string
	Delimiter string
}

type GetOnlineFeatureValuesOpt struct {
	FeatureNames []string
	EntityKeys   []string
}
