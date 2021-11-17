package metadata

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
	GroupID     int
	DBValueType string
	Description string
	ValueType   string
}

type CreateFeatureGroupOpt struct {
	Name        string
	EntityID    int
	Description string
	Category    string
}

type CreateRevisionOpt struct {
	Revision    int64
	GroupID     int
	DataTable   *string
	Anchored    bool
	Description string
}

// Update
type UpdateEntityOpt struct {
	EntityID       int
	NewDescription string
}

type UpdateFeatureOpt struct {
	FeatureID      int
	NewDescription string
}

type UpdateFeatureGroupOpt struct {
	GroupID             int
	NewDescription      *string
	NewOnlineRevisionID *int
}

type UpdateRevisionOpt struct {
	RevisionID  int
	NewRevision *int64
	NewAnchored *bool
}

// Get
type GetRevisionOpt struct {
	GroupID    *int
	Revision   *int64
	RevisionID *int
}

// List
type ListRevisionOpt struct {
	GroupID    *int
	DataTables []string
}

type ListFeatureOpt struct {
	EntityID     *int
	GroupID      *int
	FeatureIDs   *[]int
	FeatureNames *[]string
}
