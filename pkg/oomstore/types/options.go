package types

import (
	"io"
)

type CreateFeatureOpt struct {
	FeatureName string
	GroupID     int
	DBValueType string
	Description string
}

type ListFeatureOpt struct {
	EntityName   *string
	GroupName    *string
	FeatureNames *[]string
}

type UpdateFeatureOpt struct {
	FeatureName    string
	NewDescription string
}

type CreateEntityOpt struct {
	EntityName  string
	Length      int
	Description string
}

type CreateFeatureGroupOpt struct {
	GroupName   string
	EntityName  string
	Description string
}

type ExportFeatureValuesOpt struct {
	RevisionID   int
	FeatureNames []string
	Limit        *uint64
}

type ImportOpt struct {
	GroupID     int
	Description string
	DataSource  CsvDataSource
	Revision    *int64
}

type CsvDataSource struct {
	Reader    io.Reader
	Delimiter string
}

type OnlineGetOpt struct {
	FeatureNames []string
	EntityKey    string
}

type OnlineMultiGetOpt struct {
	FeatureIDs []int
	EntityKeys []string
}

type JoinOpt struct {
	FeatureIDs []int
	EntityRows <-chan EntityRow
}

type UpdateEntityOpt struct {
	EntityName     string
	NewDescription string
}

type UpdateFeatureGroupOpt struct {
	GroupName           string
	NewDescription      *string
	NewOnlineRevisionID *int
}

type SyncOpt struct {
	RevisionID int
}

type GetRevisionOpt struct {
	GroupName  *string
	Revision   *int64
	RevisionID *int
}
