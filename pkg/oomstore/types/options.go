package types

type CreateFeatureOpt struct {
	FeatureName string
	GroupName   string
	ValueType   ValueType
	Description string
}

type ListFeatureOpt struct {
	EntityNames  *[]string
	GroupNames   *[]string
	FeatureNames *[]string
}

type UpdateFeatureOpt struct {
	FeatureName    string
	NewDescription *string
}

type CreateEntityOpt struct {
	EntityName  string
	Description string
}

type ListEntityOpt struct {
	EntityNames *[]string
}

type CreateGroupOpt struct {
	GroupName        string
	EntityName       string
	Category         Category
	SnapshotInterval int
	Description      string
}

type ListGroupOpt struct {
	EntityNames *[]string
	GroupNames  *[]string
}

type ChannelExportOpt struct {
	FeatureNames []string
	UnixMilli    int64
	Limit        *uint64
}

type ExportOpt struct {
	FeatureNames   []string
	UnixMilli      int64
	Limit          *uint64
	OutputFilePath string
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
	JoinFeatureNames    []string
	ExistedFeatureNames []string
	EntityRows          <-chan EntityRow
}

type JoinOpt struct {
	FeatureNames  []string
	InputFilePath string
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
	GroupName  string
	RevisionID *int
	PurgeDelay int
}

type PushOpt struct {
	EntityKey string
	GroupName string

	// feature names without group prefix
	FeatureValues map[string]interface{}
}
