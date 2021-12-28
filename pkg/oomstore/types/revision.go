package types

import (
	"time"
)

type Revision struct {
	ID            int    `db:"id"`
	Revision      int64  `db:"revision"`
	SnapshotTable string `db:"snapshot_table"`
	CdcTable      string `db:"cdc_table"`
	Anchored      bool   `db:"anchored"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`

	GroupID int `db:"group_id"`
	Group   *Group
}

func (r *Revision) Copy() *Revision {
	if r == nil {
		return nil
	}
	copied := *r

	if copied.Group != nil {
		copied.Group = copied.Group.Copy()
	}
	return &copied
}

type RevisionList []*Revision

func (l RevisionList) Copy() RevisionList {
	if len(l) == 0 {
		return nil
	}
	copied := make(RevisionList, 0, len(l))
	for _, x := range l {
		copied = append(copied, x.Copy())
	}
	return copied
}

func (l RevisionList) Find(find func(*Revision) bool) *Revision {
	for _, r := range l {
		if find(r) {
			return r
		}
	}
	return nil
}

func (l RevisionList) Filter(filter func(*Revision) bool) (rs RevisionList) {
	for _, r := range l {
		if filter(r) {
			rs = append(rs, r)
		}
	}
	return
}

func (l RevisionList) GroupIDs() []int {
	groupIDMap := make(map[int]struct{})
	for _, r := range l {
		groupIDMap[r.GroupID] = struct{}{}
	}
	groupIDs := make([]int, 0, len(groupIDMap))
	for id := range groupIDMap {
		groupIDs = append(groupIDs, id)
	}
	return groupIDs
}
