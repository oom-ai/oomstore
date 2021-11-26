package types

import (
	"fmt"
	"strings"
	"time"
)

type Group struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Category string `db:"category"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`

	EntityID         int  `db:"entity_id"`
	OnlineRevisionID *int `db:"online_revision_id"`

	Entity *Entity
}

func (fg *Group) Copy() *Group {
	if fg == nil {
		return nil
	}
	copied := *fg

	if copied.OnlineRevisionID != nil {
		id := *copied.OnlineRevisionID
		copied.OnlineRevisionID = &id
	}
	copied.Entity = copied.Entity.Copy()

	return &copied
}

type GroupList []*Group

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

func (fg *Group) String() string {
	onlineRevisionID := "<NULL>"

	if fg.OnlineRevisionID != nil {
		onlineRevisionID = fmt.Sprint(*fg.OnlineRevisionID)
	}
	return strings.Join([]string{
		fmt.Sprintf("Name:             %s", fg.Name),
		fmt.Sprintf("Entity:           %s", fg.Entity.Name),
		fmt.Sprintf("Description:      %s", fg.Description),
		fmt.Sprintf("OnlineRevisionID: %s", onlineRevisionID),
		fmt.Sprintf("CreateTime:       %s", fg.CreateTime.Format(time.RFC3339)),
		fmt.Sprintf("ModifyTime:       %s", fg.ModifyTime.Format(time.RFC3339)),
	}, "\n")
}
