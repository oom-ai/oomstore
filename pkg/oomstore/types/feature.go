package types

import (
	"fmt"
	"strings"
	"time"
)

type Feature struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	ValueType   ValueType `db:"value_type"`
	DBValueType string    `db:"db_value_type"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`

	GroupID int `db:"group_id"`
	Group   *Group
}

func (f *Feature) Copy() *Feature {
	if f == nil {
		return nil
	}
	copied := *f

	if copied.Group != nil {
		copied.Group = copied.Group.Copy()
	}
	return &copied
}

func (f *Feature) Entity() *Entity {
	return f.Group.Entity
}

func (f *Feature) OnlineRevisionID() *int {
	return f.Group.OnlineRevisionID
}

type FeatureList []*Feature

func (l FeatureList) Copy() FeatureList {
	if len(l) == 0 {
		return nil
	}
	copied := make(FeatureList, 0, len(l))
	for _, x := range l {
		copied = append(copied, x.Copy())
	}
	return copied
}

func (l *FeatureList) Len() int { return len(*l) }

func (l *FeatureList) Names() (names []string) {
	for _, f := range *l {
		names = append(names, f.Name)
	}
	return
}

func (l *FeatureList) IDs() (ids []int) {
	for _, f := range *l {
		ids = append(ids, f.ID)
	}
	return
}

func (l FeatureList) Filter(filter func(*Feature) bool) (rs FeatureList) {
	for _, f := range l {
		if filter(f) {
			rs = append(rs, f)
		}
	}
	return
}

func (l FeatureList) Find(find func(*Feature) bool) *Feature {
	for _, f := range l {
		if find(f) {
			return f
		}
	}
	return nil
}

func (l FeatureList) GroupIDs() (ids []int) {
	for _, f := range l {
		ids = append(ids, f.GroupID)
	}
	return ids
}

func (f *Feature) String() string {
	onlineRevisionID := "<NULL>"

	if f.OnlineRevisionID() != nil {
		onlineRevisionID = fmt.Sprint(*f.OnlineRevisionID())
	}

	return strings.Join([]string{
		fmt.Sprintf("Name:             %s", f.Name),
		fmt.Sprintf("Group:            %s", f.Group.Name),
		fmt.Sprintf("Entity:           %s", f.Entity().Name),
		fmt.Sprintf("Category:         %s", f.Group.Category),
		fmt.Sprintf("DBValueType:      %s", f.DBValueType),
		fmt.Sprintf("ValueType:        %s", f.ValueType),
		fmt.Sprintf("Description:      %s", f.Description),
		fmt.Sprintf("OnlineRevisionID: %s", onlineRevisionID),
		fmt.Sprintf("CreateTime:       %s", f.CreateTime.Format(time.RFC3339)),
		fmt.Sprintf("ModifyTime:       %s", f.ModifyTime.Format(time.RFC3339)),
	}, "\n")
}
