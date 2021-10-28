package types

import (
	"fmt"
	"strings"
	"time"
)

type Feature struct {
	ID          int16  `db:"id"`
	Name        string `db:"name"`
	GroupName   string `db:"group_name"`
	ValueType   string `db:"value_type"`
	DBValueType string `db:"db_value_type"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`
}

type RichFeature struct {
	Feature
	EntityName string  `db:"entity_name"`
	Category   string  `db:"category"`
	Revision   *int64  `db:"revision"`
	DataTable  *string `db:"data_table"`
}

type FeatureList []*Feature
type RichFeatureList []*RichFeature

func (l *FeatureList) Len() int     { return len(*l) }
func (l *RichFeatureList) Len() int { return len(*l) }

func (l *FeatureList) Names() (names []string) {
	for _, f := range *l {
		names = append(names, f.Name)
	}
	return
}
func (l *RichFeatureList) Names() (names []string) {
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
func (l *RichFeatureList) Ids() (ids []int16) {
	for _, f := range *l {
		ids = append(ids, f.ID)
	}
	return
}

func (l *FeatureList) Filter(filter func(f *Feature) bool) (rs FeatureList) {
	for _, f := range *l {
		if filter(f) {
			rs = append(rs, f)
		}
	}
	return
}
func (l *RichFeatureList) Filter(filter func(f *RichFeature) bool) (rs RichFeatureList) {
	for _, f := range *l {
		if filter(f) {
			rs = append(rs, f)
		}
	}
	return
}

func (l *RichFeatureList) ToFeatureList() (rs FeatureList) {
	for _, f := range *l {
		rs = append(rs, f.AsFeature())
	}
	return
}

func (rf *RichFeature) AsFeature() *Feature {
	if rf == nil {
		return nil
	}
	return &rf.Feature
}

func (rf *RichFeature) String() string {
	var revision, dataTable string

	if rf.Revision == nil {
		revision = "<NULL>"
	} else {
		revision = fmt.Sprint(*rf.Revision)
	}

	if rf.DataTable == nil {
		dataTable = "<NULL>"
	} else {
		dataTable = *rf.DataTable
	}

	return strings.Join([]string{
		fmt.Sprintf("Name:          %s", rf.Name),
		fmt.Sprintf("Group:         %s", rf.GroupName),
		fmt.Sprintf("Entity:        %s", rf.EntityName),
		fmt.Sprintf("Category:      %s", rf.Category),
		fmt.Sprintf("DBValueType:   %s", rf.DBValueType),
		fmt.Sprintf("ValueType:     %s", rf.ValueType),
		fmt.Sprintf("Description:   %s", rf.Description),
		fmt.Sprintf("Revision:      %s", revision),
		fmt.Sprintf("DataTable:     %s", dataTable),
		fmt.Sprintf("CreateTime:    %s", rf.CreateTime.Format(time.RFC3339)),
		fmt.Sprintf("ModifyTime:    %s", rf.ModifyTime.Format(time.RFC3339)),
	}, "\n")
}
