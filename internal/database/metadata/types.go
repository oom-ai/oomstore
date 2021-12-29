package metadata

import "github.com/oom-ai/oomstore/pkg/oomstore/types"

// Create
type CreateEntityOpt struct {
	types.CreateEntityOpt
}

type CreateFeatureOpt struct {
	FeatureName string
	FullName    string
	GroupID     int
	Description string
	ValueType   types.ValueType
}

type CreateGroupOpt struct {
	GroupName   string
	EntityID    int
	Description string
	Category    types.Category
}

type CreateRevisionOpt struct {
	Revision      int64
	GroupID       int
	SnapshotTable *string
	Anchored      bool
	Description   string
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
	RevisionID       int
	NewRevision      *int64
	NewAnchored      *bool
	NewSnapshotTable *string
}

type ListFeatureOpt struct {
	EntityID         *int
	GroupID          *int
	FeatureIDs       *[]int
	FeatureFullNames *[]string
}
