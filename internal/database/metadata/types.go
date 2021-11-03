package metadata

type CreateFeatureOpt struct {
	FeatureName string
	GroupId     int16
	DBValueType string
	Description string
	ValueType   string
}

type CreateFeatureGroupOpt struct {
	Name        string
	EntityId    int16
	Description string
	Category    string
}

type CreateRevisionOpt struct {
	Revision    int64
	GroupId     int16
	GroupName   string
	DataTable   string
	Description string
}

type GetRevisionOpt struct {
	GroupName  *string
	Revision   *int64
	RevisionId *int32
}

type ListRevisionOpt struct {
	GroupName  *string
	DataTables []string
}
