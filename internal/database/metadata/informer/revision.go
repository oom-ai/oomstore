package informer

import (
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type RevisionCache struct {
	types.RevisionList
}

func (c *RevisionCache) Enrich(groupCache *GroupCache) {
	for _, r := range c.RevisionList {
		r.Group = groupCache.Find(func(g *types.FeatureGroup) bool {
			return g.ID == r.GroupID
		})
	}
}

func (c *RevisionCache) List(groupID *int) types.RevisionList {
	if groupID == nil {
		return c.RevisionList
	}
	return c.RevisionList.Filter(func(r *types.Revision) bool {
		return r.GroupID == *groupID
	})
}
