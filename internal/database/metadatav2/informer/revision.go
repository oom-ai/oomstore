package informer

import (
	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

type RevisionCache struct {
	typesv2.RevisionList
}

func (c *RevisionCache) Enrich(groupCache *GroupCache) {
	for _, r := range c.RevisionList {
		r.Group = groupCache.Find(func(g *typesv2.FeatureGroup) bool {
			return g.ID == r.GroupID
		})
	}
}

func (c *RevisionCache) List(opt metadatav2.ListRevisionOpt) typesv2.RevisionList {
	var revisions typesv2.RevisionList
	if opt.DataTables != nil {
		for _, table := range opt.DataTables {
			if r := c.Find(func(r *typesv2.Revision) bool {
				return r.DataTable == table
			}); r != nil {
				revisions = append(revisions, r)
			}
		}
	} else {
		revisions = c.RevisionList
	}

	if opt.GroupID != nil {
		revisions = revisions.Filter(func(r *typesv2.Revision) bool {
			return r.GroupID == *opt.GroupID
		})
	}
	return revisions
}

func (c *RevisionCache) GetGroup(groupID int16) typesv2.RevisionList {
	return c.Filter(func(r *typesv2.Revision) bool {
		return r.GroupID == groupID
	})
}
