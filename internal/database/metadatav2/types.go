package metadatav2

type RevisionRange struct {
	MinRevision int64  `db:"min_revision"`
	MaxRevision int64  `db:"max_revision"`
	DataTable   string `db:"data_table"`
}

// Create
type CreateEntityOpt struct {
	Name        string
	Length      int
	Description string
}

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

type CreateRevisionOpt struct {
	Revision    int64
	GroupId     int16
	DataTable   string
	Anchored    bool
	Description string
}

// Update
type UpdateEntityOpt struct {
	EntityID       int16
	NewDescription string
}

type UpdateFeatureOpt struct {
	FeatureID      int16
	NewDescription string
}

type UpdateFeatureGroupOpt struct {
	GroupID          int16
	Description      *string
	OnlineRevisionId *int32
}

type UpdateRevisionOpt struct {
	RevisionID  int32
	NewRevision *int64
	NewAnchored *bool
}

// Get
type GetRevisionOpt struct {
	GroupName  *string
	Revision   *int64
	RevisionId *int32
}

// List
type ListRevisionOpt struct {
	GroupName  *string
	DataTables []string
}

type ListFeatureOpt struct {
	EntityName   *string
	GroupName    *string
	FeatureNames []string
}
