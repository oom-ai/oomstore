package informer

import (
	"sort"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type FeatureCache struct {
	types.FeatureList
	nameIdx map[string]*types.Feature
}

func NewFeatureCache(features types.FeatureList) *FeatureCache {
	sort.Slice(features, func(i, j int) bool {
		return features[i].ID < features[j].ID
	})
	return &FeatureCache{FeatureList: features, nameIdx: nil}
}

func (c *FeatureCache) Enrich(groupCache *GroupCache) {
	nameIdx := make(map[string]*types.Feature)

	for _, f := range c.FeatureList {
		f.Group = groupCache.Find(func(g *types.Group) bool {
			return g.ID == f.GroupID
		})
		nameIdx[f.FullName()] = f
	}

	c.nameIdx = nameIdx
}

func (c *FeatureCache) Get(featureID int) *types.Feature {
	l := c.FeatureList
	pos := sort.Search(len(l), func(i int) bool {
		return l[i].ID == featureID
	})
	if pos >= 0 && pos < len(l) {
		return l[pos]
	}
	return nil
}

func (c *FeatureCache) List(opt metadata.ListCachedFeatureOpt) types.FeatureList {
	features := c.FeatureList

	// filter featureNames
	if opt.FullNames != nil {
		var tmp types.FeatureList
		for _, name := range *opt.FullNames {
			if f, ok := c.nameIdx[name]; ok {
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
