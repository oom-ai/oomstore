package metadatav2

type CreateFeatureOpt struct {
	FeatureName string
	GroupID     int16
	DBValueType string
	Description string
	ValueType   string
}

type CreateFeatureGroupOpt struct {
	Name        string
	EntityID    int16
	Description string
	Category    string
}

type UpdateFeatureGroupOpt struct {
	GroupID          int16
	Description      *string
	OnlineRevisionId *int32
}

type CreateRevisionOpt struct {
	Revision    int64
	GroupId     int16
	DataTable   string
	Anchored    bool
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

type UpdateRevisionOpt struct {
	RevisionID  int32
	NewRevision *int64
	NewAnchored *bool
}
