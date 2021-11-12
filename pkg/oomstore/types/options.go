package types

import (
	"io"
)

type CreateFeatureOpt struct {
	FeatureName string
	GroupID     string
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
	GroupID      int16
	Revision     *int64
	FeatureNames []string
	Limit        *uint64
}

type ImportBatchFeaturesOpt struct {
	GroupID     int16
	Description string
	DataSource  CsvDataSource
	Revision    *int64
}

type CsvDataSource struct {
	Reader    io.Reader
	Delimiter string
}

type GetOnlineFeatureValuesOpt struct {
	FeatureIDs []int16
	EntityKey  string
}

type MultiGetOnlineFeatureValuesOpt struct {
	FeatureIDs []int16
	EntityKeys []string
}

type GetHistoricalFeatureValuesOpt struct {
	FeatureIDs []int16
	EntityRows <-chan EntityRow
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

type SyncOpt struct {
	RevisionId int32
}

type GetRevisionOpt struct {
	GroupName  *string
	Revision   *int64
	RevisionId *int32
}
