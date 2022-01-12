package types

type CreateFeatureOpt struct {
	FeatureName string
	GroupName   string
	ValueType   ValueType
	Description string
}

type ListFeatureOpt struct {
	EntityName       *string
	GroupName        *string
	FeatureFullNames *[]string
}

type UpdateFeatureOpt struct {
	FeatureFullName string
	NewDescription  *string
}

type CreateEntityOpt struct {
	EntityName  string
	Description string
}

type CreateGroupOpt struct {
	GroupName   string
	EntityName  string
	Category    Category
	Description string
}

type ChannelExportBatchOpt struct {
	RevisionID   int
	FeatureNames []string
	Limit        *uint64
}

type ChannelExportStreamOpt struct {
	UnixMilli        int64
	FeatureFullNames []string
	Limit            *uint64
}

type ExportBatchOpt struct {
	RevisionID     int
	FeatureNames   []string
	Limit          *uint64
	OutputFilePath string
}

type OnlineGetOpt struct {
	FeatureFullNames []string
	EntityKey        string
}

type OnlineMultiGetOpt struct {
	FeatureFullNames []string
	EntityKeys       []string
}

type ChannelJoinOpt struct {
	FeatureFullNames []string
	EntityRows       <-chan EntityRow
	ValueNames       []string
}

type JoinOpt struct {
	FeatureFullNames []string
	InputFilePath    string
	OutputFilePath   string
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

type PushOpt struct {
	EntityKey     string
	GroupName     string
	FeatureNames  []string
	FeatureValues []interface{}
}
