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
	Name        string
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
	GroupID     int16
	DataTable   *string
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
	GroupID             int16
	NewDescription      *string
	NewOnlineRevisionID *int32
}

type UpdateRevisionOpt struct {
	RevisionID  int32
	NewRevision *int64
	NewAnchored *bool
}

// Get
type GetRevisionOpt struct {
	GroupID    *int16
	Revision   *int64
	RevisionID *int32
}

// List
type ListRevisionOpt struct {
	GroupID    *int16
	DataTables []string
}

type ListFeatureOpt struct {
	EntityID     *int16
	GroupID      *int16
	FeatureIDs   *[]int16
	FeatureNames *[]string
}
