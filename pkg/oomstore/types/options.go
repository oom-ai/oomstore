package types

import (
	"io"
)

type CreateFeatureOpt struct {
	FeatureName string
	GroupName   string
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
	NewDescription *string
}

type CreateEntityOpt struct {
	EntityName  string
	Length      int
	Description string
}

type CreateGroupOpt struct {
	GroupName   string
	EntityName  string
	Description string
}

type ChannelExportOpt struct {
	RevisionID   int
	FeatureNames []string
	Limit        *uint64
}

type ExportOpt struct {
	RevisionID     int
	FeatureNames   []string
	Limit          *uint64
	OutputFilePath string
}

type ChannelImport struct {
	GroupName   string
	Description string
	DataSource  CsvDataSource
	Revision    *int64
}

type ImportOpt struct {
	GroupName   string
	Description string
	DataSource  interface{}
	Revision    *int64
}

type CsvDataSource struct {
	Reader    io.Reader
	Delimiter string
}

type CsvDataSourceWithFile struct {
	InputFilePath string
	Delimiter     string
}

type OnlineGetOpt struct {
	FeatureNames []string
	EntityKey    string
}

type OnlineMultiGetOpt struct {
	FeatureNames []string
	EntityKeys   []string
}

type ChannelJoinOpt struct {
	FeatureNames []string
	EntityRows   <-chan EntityRow
	ValueNames   []string
}

type JoinOpt struct {
	FeatureNames   []string
	InputFilePath  string
	OutputFilePath string
}

type UpdateEntityOpt struct {
	EntityName     string
	NewDescription *string
}

type UpdateGroupOpt struct {
	GroupName           string
	NewDescription      *string
	NewOnlineRevisionID *int
}

type SyncOpt struct {
	RevisionID int
	PurgeDelay int
}
