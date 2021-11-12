package typesv2

import (
	"fmt"
	"strings"
	"time"
)

type Feature struct {
	ID          int16  `db:"id"`
	Name        string `db:"name"`
	ValueType   string `db:"value_type"`
	DBValueType string `db:"db_value_type"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`

	GroupID int16 `db:"group_id"`
	Group   *FeatureGroup
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

func (f *Feature) OnlineRevision() *Revision {
	return f.Group.OnlineRevision
}

type FeatureList []*Feature

func (l *FeatureList) Len() int { return len(*l) }

func (l *FeatureList) Names() (names []string) {
	for _, f := range *l {
		names = append(names, f.Name)
	}
	return
}

func (l *FeatureList) Ids() (ids []int16) {
	for _, f := range *l {
		ids = append(ids, f.ID)
	}
	return
}

func (l *FeatureList) Filter(filter func(*Feature) bool) (rs FeatureList) {
	for _, f := range *l {
		if filter(f) {
			rs = append(rs, f)
		}
	}
	return
}

func (l *FeatureList) Find(find func(*Feature) bool) *Feature {
	for _, f := range *l {
		if find(f) {
			return f
		}
	}
	return nil
}

func (f *Feature) String() string {
	onlineRevision := "<NULL>"

	if f.OnlineRevision() != nil {
		onlineRevision = fmt.Sprint(*f.OnlineRevision())
	}

	return strings.Join([]string{
		fmt.Sprintf("Name:            %s", f.Name),
		fmt.Sprintf("Group:           %s", f.Group.Name),
		fmt.Sprintf("Entity:          %s", f.Entity().Name),
		fmt.Sprintf("Category:        %s", f.Group.Category),
		fmt.Sprintf("DBValueType:     %s", f.DBValueType),
		fmt.Sprintf("ValueType:       %s", f.ValueType),
		fmt.Sprintf("Description:     %s", f.Description),
		fmt.Sprintf("Online Revision: %s", onlineRevision),
		fmt.Sprintf("CreateTime:      %s", f.CreateTime.Format(time.RFC3339)),
		fmt.Sprintf("ModifyTime:      %s", f.ModifyTime.Format(time.RFC3339)),
	}, "\n")
}

func FeatureCsvHeader() string {
	return strings.Join([]string{
		"Name",
		"Group",
		"Entity",
		"Category",
		"DBValueType",
		"ValueType",
		"Description",
		"OnlineRevision",
		"CreateTime",
		"ModifyTime",
	}, ",")
}

func (f *Feature) ToCsvRecord() string {
	onlineRevision := "<NULL>"

	if r := f.OnlineRevision(); r != nil {
		onlineRevision = fmt.Sprint(r.Revision)
	}

	return strings.Join([]string{
		f.Name,
		f.Group.Name,
		f.Entity().Name,
		f.Group.Category,
		f.DBValueType,
		f.ValueType,
		f.Description,
		onlineRevision,
		f.CreateTime.Format(time.RFC3339),
		f.ModifyTime.Format(time.RFC3339),
	}, ",")
}
