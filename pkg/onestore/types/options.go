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
	Length      int
	Description string
}

type CreateFeatureGroupOpt struct {
	Name        string
	EntityName  string
	Description string
}

type WalkFeatureValuesFunc = func(header []string, key string, values []interface{}) error

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
	Delimiter string
}

type GetOnlineFeatureValuesOpt struct {
	FeatureNames []string
	EntityKey    string
}

type GetOnlineFeatureValuesWithMultiEntityKeysOpt struct {
	FeatureNames []string
	EntityKeys   []string
}

type EntityRow struct {
	EntityKey string `db:"entity_key"`
	UnixTime  int64  `db:"unix_time"`
}

type GetHistoricalFeatureValuesOpt struct {
	FeatureNames []string
	EntityRows   []EntityRow
}

type UpdateEntityOpt struct {
	EntityName     string
	NewDescription string
}

type UpdateFeatureGroupOpt struct {
	GroupName      string
	NewDescription string
}
