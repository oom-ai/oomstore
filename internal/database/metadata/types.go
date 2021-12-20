package metadata

import "github.com/oom-ai/oomstore/pkg/oomstore/types"

type RevisionRange struct {
	MinRevision int64  `db:"min_revision"`
	MaxRevision int64  `db:"max_revision"`
	DataTable   string `db:"data_table"`
}

// Create
type CreateEntityOpt struct {
	types.CreateEntityOpt
}

type CreateFeatureOpt struct {
	FeatureName string
	GroupID     int
	DBValueType string
	Description string
	ValueType   string
}

type CreateGroupOpt struct {
	GroupName   string
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
	NewDescription *string
}

type UpdateFeatureOpt struct {
	FeatureID      int
	NewDescription *string
}

type UpdateGroupOpt struct {
	GroupID             int
	NewDescription      *string
	NewOnlineRevisionID *int
}

type UpdateRevisionOpt struct {
	RevisionID  int
	NewRevision *int64
	NewAnchored *bool
}

type ListFeatureOpt struct {
	EntityID     *int
	GroupID      *int
	FeatureIDs   *[]int
	FeatureNames *[]string
}
