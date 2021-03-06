package metadata

import (
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// Create
type CreateEntityOpt struct {
	types.CreateEntityOpt
}

type CreateFeatureOpt struct {
	FeatureName string
	GroupID     int
	Description string
	ValueType   types.ValueType
}

type CreateGroupOpt struct {
	GroupName   string
	EntityID    int
	Description string
	Category    types.Category

	SnapshotInterval int
}

type CreateRevisionOpt struct {
	Revision      int64
	GroupID       int
	SnapshotTable *string
	CdcTable      *string
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
	NewSnapshotInterval *int
	NewDescription      *string
	NewOnlineRevisionID *int
}

type UpdateRevisionOpt struct {
	RevisionID       int
	NewRevision      *int64
	NewAnchored      *bool
	NewSnapshotTable *string
	NewCdcTable      *string
}

type ListFeatureOpt struct {
	EntityIDs  *[]int
	GroupIDs   *[]int
	FeatureIDs *[]int
}

type ListGroupOpt struct {
	EntityIDs *[]int
	GroupIDs  *[]int
}

type ListCachedFeatureOpt struct {
	FullNames *[]string
	GroupName *string
	GroupID   *int
}
