package online

import (
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type GetOpt struct {
	EntityKey string
	Group     types.Group
	Features  types.FeatureList

	// Only works when get batch features, it should be nil when get stream features
	RevisionID *int
}

func (g *GetOpt) Validate() error {
	if g.Group.Category == types.CategoryBatch && g.RevisionID == nil {
		return errdefs.Errorf("invalid GetOpt: the revisionID of get batch feature cannot be null")
	}
	return nil
}

type GetByGroupOpt struct {
	EntityKey   string
	Group       types.Group
	ListFeature func(groupID int) types.FeatureList

	// Only works when get batch features, it should be nil when get stream features
	RevisionID *int
}

func (g *GetByGroupOpt) Validate() error {
	if g.Group.Category == types.CategoryBatch && g.RevisionID == nil {
		return errdefs.Errorf("invalid GetByGroupOpt: the revisionID of get batch feature cannot be null")
	}
	if g.ListFeature == nil {
		return errdefs.Errorf("invalid GetByGroupOpt: the ListCachedFeature function cannot be null")
	}
	return nil
}

type MultiGetOpt struct {
	EntityKeys []string
	Group      types.Group
	Features   types.FeatureList

	// Only works when get batch features, it should be nil when get stream features
	RevisionID *int
}

func (m *MultiGetOpt) Validate() error {
	if m.Group.Category == types.CategoryBatch && m.RevisionID == nil {
		return errdefs.Errorf("invalid MultiGetOpt: the revisionID of get batch feature cannot be null")
	}
	return nil
}

type ImportOpt struct {
	Group        types.Group
	Features     types.FeatureList
	ExportStream <-chan types.ExportRecord
	RevisionID   *int
}

type PushOpt struct {
	EntityName    string
	EntityKey     string
	GroupID       int
	Features      types.FeatureList
	FeatureValues []interface{}
}

type CreateTableOpt struct {
	EntityName string
	TableName  string
	Features   types.FeatureList
}
