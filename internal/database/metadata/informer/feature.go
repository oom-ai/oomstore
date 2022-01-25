package informer

import (
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type FeatureCache struct {
	types.FeatureList
}

func (c *FeatureCache) Enrich(groupCache *GroupCache) {
	for _, f := range c.FeatureList {
		f.Group = groupCache.Find(func(g *types.Group) bool {
			return g.ID == f.GroupID
		})
	}
}

func (c *FeatureCache) List(opt metadata.ListCachedFeatureOpt) types.FeatureList {
	features := c.FeatureList

	// filter featureNames
	if opt.FullNames != nil {
		var tmp types.FeatureList
		for _, fullName := range *opt.FullNames {
			if f := features.Find(func(f *types.Feature) bool {
				return f.FullName() == fullName
			}); f != nil {
				tmp = append(tmp, f)
			}
		}
		features = tmp
	}
	// filter groupName
	if opt.GroupName != nil {
		features = features.Filter(func(f *types.Feature) bool {
			return f.Group.Name == *opt.GroupName
		})
	}
	// filter groupID
	if opt.GroupID != nil {
		features = features.Filter(func(f *types.Feature) bool {
			return f.Group.ID == *opt.GroupID
		})
	}
	return features
}
