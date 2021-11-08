package types

import "io"

type CreateFeatureOpt struct {
	FeatureName string
	GroupName   string
	DBValueType string
	Description string
}

type ListFeatureOpt struct {
	EntityName   *string
	GroupName    *string
	FeatureNames []string
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

type ExportFeatureValuesOpt struct {
	GroupName     string
	GroupRevision *int64
	FeatureNames  []string
	Limit         *uint64
}

type ImportBatchFeaturesOpt struct {
	GroupName   string
	Description string
	DataSource  CsvDataSource
	Revision    *int64
}

type CsvDataSource struct {
	Reader    io.Reader
	Delimiter string
}

type GetOnlineFeatureValuesOpt struct {
	FeatureNames []string
	EntityKey    string
}

type MultiGetOnlineFeatureValuesOpt struct {
	FeatureNames []string
	EntityKeys   []string
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
	GroupName        string
	Description      *string
	OnlineRevisionId *int32
}

type MaterializeOpt struct {
	GroupName     string
	GroupRevision *int64
}
