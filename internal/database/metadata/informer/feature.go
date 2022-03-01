package informer

import (
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type FeatureCache struct {
	types.FeatureList
	nameIdx map[string]*types.Feature
	idIdx   map[int]*types.Feature
}

func NewFeatureCache(features types.FeatureList) *FeatureCache {
	return &FeatureCache{FeatureList: features, nameIdx: nil, idIdx: nil}
}

func (c *FeatureCache) Enrich(groupCache *GroupCache) {
	nameIdx := make(map[string]*types.Feature)
	idIdx := make(map[int]*types.Feature)

	for _, f := range c.FeatureList {
		f.Group = groupCache.Find(func(g *types.Group) bool {
			return g.ID == f.GroupID
		})
		nameIdx[f.FullName()] = f
		idIdx[f.ID] = f
	}

	c.nameIdx = nameIdx
	c.idIdx = idIdx
}

func (c *FeatureCache) Get(featureID int) *types.Feature {
	return c.idIdx[featureID]
}

func (c *FeatureCache) GetByName(fullName string) *types.Feature {
	return c.nameIdx[fullName]
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
