package types

import (
	"fmt"
	"strings"
	"time"
)

type FeatureGroup struct {
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

func (fg *FeatureGroup) Copy() *FeatureGroup {
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

type FeatureGroupList []*FeatureGroup

func (l FeatureGroupList) Copy() FeatureGroupList {
	if len(l) == 0 {
		return nil
	}
	copied := make(FeatureGroupList, 0, len(l))
	for _, x := range l {
		copied = append(copied, x.Copy())
	}
	return copied
}

func (l *FeatureGroupList) Find(find func(*FeatureGroup) bool) *FeatureGroup {
	for _, g := range *l {
		if find(g) {
			return g
		}
	}
	return nil
}

func (l *FeatureGroupList) Filter(filter func(*FeatureGroup) bool) (rs FeatureGroupList) {
	for _, g := range *l {
		if filter(g) {
			rs = append(rs, g)
		}
	}
	return
}

func (fg *FeatureGroup) String() string {
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
