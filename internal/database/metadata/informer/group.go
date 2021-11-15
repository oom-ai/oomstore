package informer

import (
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type GroupCache struct {
	types.FeatureGroupList
}

func (c *GroupCache) Enrich(entityCache *EntityCache) {
	for _, g := range c.FeatureGroupList {
		g.Entity = entityCache.Find(func(e *types.Entity) bool {
			return e.ID == g.EntityID
		})
	}
}

func (c *GroupCache) List(entityID *int16) types.FeatureGroupList {
	if entityID == nil {
		return c.FeatureGroupList
	}
	return c.FeatureGroupList.Filter(func(g *types.FeatureGroup) bool {
		return g.Entity.ID == *entityID
	})
}
