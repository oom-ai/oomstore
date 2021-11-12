package typesv2

import (
	"time"
)

type Revision struct {
	ID        int32  `db:"id"`
	Revision  int64  `db:"revision"`
	DataTable string `db:"data_table"`
	Anchored  bool   `db:"anchored"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`

	GroupID int16 `db:"group_id"`
	Group   *FeatureGroup
}

func (r *Revision) Copy() *Revision {
	return r.copyWith(nil)
}

func (r *Revision) copyWith(group *FeatureGroup) *Revision {
	if r == nil {
		return nil
	}

	copied := *r
	if group != nil {
		copied.Group = group
	} else if copied.Group != nil {
		copied.Group = r.Group.copyWith(&copied)
	}
	return &copied
}

type RevisionList []*Revision

func (l *RevisionList) Find(find func(*Revision) bool) *Revision {
	for _, r := range *l {
		if find(r) {
			return r
		}
	}
	return nil
}

func (l *RevisionList) Filter(filter func(*Revision) bool) (rs RevisionList) {
	for _, r := range *l {
		if filter(r) {
			rs = append(rs, r)
		}
	}
	return
}
