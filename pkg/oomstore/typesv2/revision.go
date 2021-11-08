package typesv2

import "time"

type Revision struct {
	ID        int32  `db:"id"`
	Revision  int64  `db:"revision"`
	DataTable string `db:"data_table"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`

	GroupID int16 `db:"group_id"`
	Group   *FeatureGroup
}

type RevisionList []*Revision

func (l *RevisionList) Find(find func(r *Revision) bool) *Revision {
	for _, r := range *l {
		if find(r) {
			return r
		}
	}
	return nil
}

func (l *RevisionList) Filter(filter func(r *Revision) bool) (rs RevisionList) {
	for _, r := range *l {
		if filter(r) {
			rs = append(rs, r)
		}
	}
	return
}
