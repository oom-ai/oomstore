package informer

import (
	"github.com/oom-ai/oomstore/internal/database/metadata"
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

func (c *RevisionCache) List(opt metadata.ListRevisionOpt) types.RevisionList {
	var revisions types.RevisionList
	if opt.DataTables != nil {
		for _, table := range opt.DataTables {
			if r := c.Find(func(r *types.Revision) bool {
				return r.DataTable == table
			}); r != nil {
				revisions = append(revisions, r)
			}
		}
	} else {
		revisions = c.RevisionList
	}

	if opt.GroupID != nil {
		revisions = revisions.Filter(func(r *types.Revision) bool {
			return r.GroupID == *opt.GroupID
		})
	}
	return revisions
}
