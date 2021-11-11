package informer

import (
	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

type FeatureCache struct {
	typesv2.FeatureList
}

func (c *FeatureCache) Enrich(groupCache *GroupCache) {
	for _, f := range c.FeatureList {
		f.Group = groupCache.Find(func(g *typesv2.FeatureGroup) bool {
			return g.ID == f.GroupID
		})
	}
}

func (c *FeatureCache) List(opt metadatav2.ListFeatureOpt) typesv2.FeatureList {
	var features typesv2.FeatureList

	// filter ids
	if opt.FeatureIDs != nil {
		for _, id := range opt.FeatureIDs {
			if f := c.Find(func(f *typesv2.Feature) bool {
				return f.ID == id
			}); f != nil {
				features = append(features, f)
			}
		}
	} else {
		features = c.FeatureList
	}

	// filter entity
	if opt.EntityID != nil {
		features = features.Filter(func(f *typesv2.Feature) bool {
			return f.Entity().ID == *opt.EntityID
		})
	}

	// filter group
	if opt.GroupID != nil {
		features = features.Filter(func(f *typesv2.Feature) bool {
			return f.Group.ID == *opt.GroupID
		})
	}
	return features
}
