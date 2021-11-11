package informer

import (
	"math"
	"sort"

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
			return r.Group.ID == *opt.GroupID
		})
	}
	return revisions
}

func (c *RevisionCache) GetGroup(groupID int16) typesv2.RevisionList {
	return c.Filter(func(r *typesv2.Revision) bool {
		return r.Group.ID == groupID
	})
}

func (c *RevisionCache) MaxRevision(groupID int16) *typesv2.Revision {
	revisions := c.GetGroup(groupID)
	if revisions == nil {
		return nil
	}

	var max *typesv2.Revision
	for _, r := range revisions {
		if max == nil || max.Revision < r.Revision {
			max = r
		}
	}
	return max
}

func (c *RevisionCache) BuildRevisionRanges(groupID int16) []*metadatav2.RevisionRange {
	revisionIndex := c.GetGroup(groupID)
	if len(revisionIndex) == 0 {
		return nil
	}

	var revisions typesv2.RevisionList
	revisions = append(revisions, revisionIndex...)
	sort.Slice(revisions, func(i, j int) bool {
		return revisions[i].Revision < revisions[j].Revision
	})

	var ranges []*metadatav2.RevisionRange
	for i := 1; i < len(revisions); i++ {
		ranges = append(ranges, &metadatav2.RevisionRange{
			MinRevision: revisions[i-1].Revision,
			MaxRevision: revisions[i].Revision,
			DataTable:   revisions[i-1].DataTable,
		})
	}

	return append(ranges, &metadatav2.RevisionRange{
		MinRevision: revisions[len(revisions)-1].Revision,
		MaxRevision: revisions[math.MaxInt64].Revision,
		DataTable:   revisions[len(revisions)-1].DataTable,
	})
}