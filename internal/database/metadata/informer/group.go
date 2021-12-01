package informer

import (
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type GroupCache struct {
	types.GroupList
}

func (c *GroupCache) Enrich(entityCache *EntityCache) {
	for _, g := range c.GroupList {
		g.Entity = entityCache.Find(func(e *types.Entity) bool {
			return e.ID == g.EntityID
		})
	}
}
