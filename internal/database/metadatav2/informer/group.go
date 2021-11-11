package informer

import (
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

type GroupCache struct {
	typesv2.FeatureGroupList
}

func (c *GroupCache) Enrich(entityCache *EntityCache, revisionCache *RevisionCache) {
	for _, g := range c.FeatureGroupList {
		g.Entity = entityCache.Find(func(e *typesv2.Entity) bool {
			return e.ID == g.EntityID
		})

		if g.OnlineRevisionID != nil {
			g.OnlineRevision = revisionCache.Find(func(r *typesv2.Revision) bool {
				return r.ID == *g.OnlineRevisionID
			})
		}
	}
}

func (c *GroupCache) List(entityID *int16) []*typesv2.FeatureGroup {
	if entityID == nil {
		return c.FeatureGroupList
	}
	return c.FeatureGroupList.Filter(func(g *typesv2.FeatureGroup) bool {
		return g.Entity.ID == *entityID
	})
}
