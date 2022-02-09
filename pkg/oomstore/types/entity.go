package types

import (
	"time"
)

type Entity struct {
	ID   int    `db:"id"`
	Name string `db:"name"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`
}

func (e *Entity) Copy() *Entity {
	if e == nil {
		return nil
	}
	copied := *e
	return &copied
}

type EntityList []*Entity

func (l EntityList) Copy() EntityList {
	if len(l) == 0 {
		return nil
	}
	copied := make(EntityList, 0, len(l))
	for _, x := range l {
		copied = append(copied, x.Copy())
	}
	return copied
}

func (l EntityList) Len() int {
	return len(l)
}

func (l EntityList) Find(find func(*Entity) bool) *Entity {
	for _, e := range l {
		if find(e) {
			return e
		}
	}
	return nil
}

func (l EntityList) Filter(filter func(*Entity) bool) (rs EntityList) {
	for _, e := range l {
		if filter(e) {
			rs = append(rs, e)
		}
	}
	return
}

func (l EntityList) IDs() (ids []int) {
	for _, entity := range l {
		ids = append(ids, entity.ID)
	}
	return
}

func (l EntityList) Names() (names []string) {
	for _, entity := range l {
		names = append(names, entity.Name)
	}
	return
}
