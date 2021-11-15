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
		f.Group = groupCache.Find(func(g *types.FeatureGroup) bool {
			return g.ID == f.GroupID
		})
	}
}

func (c *FeatureCache) List(opt metadata.ListFeatureOpt) types.FeatureList {
	features := c.FeatureList

	// filter ids
	if opt.FeatureIDs != nil {
		var tmp types.FeatureList
		for _, id := range *opt.FeatureIDs {
			if f := c.Find(func(f *types.Feature) bool {
				return f.ID == id
			}); f != nil {
				tmp = append(tmp, f)
			}
		}
		features = tmp
	}

	// filter names
	if opt.FeatureNames != nil {
		var tmp types.FeatureList
		for _, name := range *opt.FeatureNames {
			if f := features.Find(func(f *types.Feature) bool {
				return f.Name == name
			}); f != nil {
				tmp = append(tmp, f)
			}
		}
		features = tmp
	}

	// filter entity
	if opt.EntityID != nil {
		features = features.Filter(func(f *types.Feature) bool {
			return f.Entity().ID == *opt.EntityID
		})
	}

	// filter group
	if opt.GroupID != nil {
		features = features.Filter(func(f *types.Feature) bool {
			return f.Group.ID == *opt.GroupID
		})
	}
	return features
}
