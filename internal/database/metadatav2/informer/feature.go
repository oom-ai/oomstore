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
	features := c.FeatureList

	if opt.FeatureIDs == nil && opt.FeatureNames == nil && opt.EntityID == nil && opt.GroupID == nil {
		features = make(typesv2.FeatureList, len(c.FeatureList))
		copy(features, c.FeatureList)
		return features
	}

	// filter ids
	if opt.FeatureIDs != nil {
		var tmp typesv2.FeatureList
		for _, id := range *opt.FeatureIDs {
			if f := c.Find(func(f *typesv2.Feature) bool {
				return f.ID == id
			}); f != nil {
				tmp = append(features, f)
			}
		}
		features = tmp
	}

	// filter names
	if opt.FeatureNames != nil {
		var tmp typesv2.FeatureList
		for _, name := range *opt.FeatureNames {
			if f := features.Find(func(f *typesv2.Feature) bool {
				return f.Name == name
			}); f != nil {
				tmp = append(tmp, f)
			}
		}
		features = tmp
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
