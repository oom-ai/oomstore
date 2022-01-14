package informer

import (
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

func (c *FeatureCache) List(fullNames *[]string) types.FeatureList {
	features := c.FeatureList

	// filter names
	if fullNames != nil {
		var tmp types.FeatureList
		for _, fullName := range *fullNames {
			if f := features.Find(func(f *types.Feature) bool {
				return f.FullName() == fullName
			}); f != nil {
				tmp = append(tmp, f)
			}
		}
		features = tmp
	}

	return features
}
