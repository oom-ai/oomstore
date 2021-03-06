package types

import (
	"time"
)

type Group struct {
	ID       int      `db:"id"`
	Name     string   `db:"name"`
	Category Category `db:"category"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`

	SnapshotInterval int `db:"snapshot_interval"`

	EntityID         int  `db:"entity_id"`
	OnlineRevisionID *int `db:"online_revision_id"`

	Entity *Entity
}

func (g *Group) Copy() *Group {
	if g == nil {
		return nil
	}
	copied := *g

	if copied.OnlineRevisionID != nil {
		id := *copied.OnlineRevisionID
		copied.OnlineRevisionID = &id
	}
	copied.Entity = copied.Entity.Copy()

	return &copied
}

type GroupList []*Group

func (l GroupList) Len() int {
	return len(l)
}

func (l GroupList) Copy() GroupList {
	if len(l) == 0 {
		return nil
	}
	copied := make(GroupList, 0, len(l))
	for _, x := range l {
		copied = append(copied, x.Copy())
	}
	return copied
}

func (l GroupList) Find(find func(*Group) bool) *Group {
	for _, g := range l {
		if find(g) {
			return g
		}
	}
	return nil
}

func (l GroupList) Filter(filter func(*Group) bool) (rs GroupList) {
	for _, g := range l {
		if filter(g) {
			rs = append(rs, g)
		}
	}
	return
}

func (l GroupList) EntityIDs() []int {
	entityIDMap := make(map[int]struct{})
	for _, g := range l {
		entityIDMap[g.EntityID] = struct{}{}
	}
	entityIDs := make([]int, 0, len(entityIDMap))
	for id := range entityIDMap {
		entityIDs = append(entityIDs, id)
	}
	return entityIDs
}

func (l GroupList) IDs() (ids []int) {
	for _, group := range l {
		ids = append(ids, group.ID)
	}
	return
}

func (l GroupList) Names() (names []string) {
	for _, group := range l {
		names = append(names, group.Name)
	}
	return
}
