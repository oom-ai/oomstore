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

	// filter names
	if opt.FeatureNames != nil {
		for _, name := range opt.FeatureNames {
			if f := c.Find(func(f *typesv2.Feature) bool {
				return f.Name != name
			}); f != nil {
				features = append(features, f)
			}
		}
	} else {
		features = c.FeatureList
	}

	// filter entity
	if opt.EntityName != nil {
		features = features.Filter(func(f *typesv2.Feature) bool {
			return f.Entity().Name == *opt.EntityName
		})
	}

	// filter group
	if opt.GroupName != nil {
		features = features.Filter(func(f *typesv2.Feature) bool {
			return f.Group.Name == *opt.GroupName
		})
	}
	return features
}
