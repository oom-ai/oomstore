package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	BatchFeatureCategory  = "batch"
	StreamFeatureCategory = "stream"
)

type Entity struct {
	ID     int16  `db:"id"`
	Name   string `db:"name"`
	Length int    `db:"length"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`
}

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

type FeatureGroup struct {
	ID               int16   `db:"id"`
	Name             string  `db:"name"`
	EntityName       string  `db:"entity_name"`
	Revision         *int64  `db:"revision"`
	OnlineRevisionID *int32  `db:"online_revision_id"`
	Category         string  `db:"category"`
	DataTable        *string `db:"data_table"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`
}

type Revision struct {
	ID        int32  `db:"id"`
	Revision  int64  `db:"revision"`
	GroupName string `db:"group_name"`
	DataTable string `db:"data_table"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`
}

type FeatureList []*Feature
type RichFeatureList []*RichFeature

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

func (l *FeatureList) Filter(filter func(f *Feature) bool) (rs FeatureList) {
	for _, f := range *l {
		if filter(f) {
			rs = append(rs, f)
		}
	}
	return
}

func (l *RichFeatureList) Names() (names []string) {
	for _, f := range *l {
		names = append(names, f.Name)
	}
	return
}

func (l *RichFeatureList) Ids() (ids []int16) {
	for _, f := range *l {
		ids = append(ids, f.ID)
	}
	return
}

func (l *RichFeatureList) ToFeatureList() (rs FeatureList) {
	for _, f := range *l {
		rs = append(rs, f.ToFeature())
	}
	return
}

func (rf *RichFeature) ToFeature() *Feature {
	if rf == nil {
		return nil
	}
	return &rf.Feature
}

func (l *RichFeatureList) Filter(filter func(f *RichFeature) bool) (rs RichFeatureList) {
	for _, f := range *l {
		if filter(f) {
			rs = append(rs, f)
		}
	}
	return
}

type FeatureKV struct {
	FeatureName  string
	FeatureValue interface{}
}

func NewFeatureKV(name string, value interface{}) FeatureKV {
	return FeatureKV{
		FeatureName:  name,
		FeatureValue: value,
	}
}

type FeatureValueMap map[string]interface{}

type FeatureDataSet map[string][]FeatureKV

func NewFeatureDataSet() FeatureDataSet {
	return make(map[string][]FeatureKV)
}

type EntityRowWithFeatures struct {
	EntityRow
	FeatureValues []FeatureKV
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

func (e *Entity) String() string {
	return strings.Join([]string{
		fmt.Sprintf("Name:          %s", e.Name),
		fmt.Sprintf("Length:        %d", e.Length),
		fmt.Sprintf("Description:   %s", e.Description),
		fmt.Sprintf("CreateTime:    %s", e.CreateTime.Format(time.RFC3339)),
		fmt.Sprintf("ModifyTime:    %s", e.ModifyTime.Format(time.RFC3339)),
	}, "\n")
}

func (fg *FeatureGroup) String() string {
	revision := "NULL"
	dataTable := "NULL"
	if fg.Revision != nil {
		revision = strconv.Itoa(int(*fg.Revision))
	}
	if fg.DataTable != nil {
		dataTable = *fg.DataTable
	}
	return strings.Join([]string{
		fmt.Sprintf("Name:          %s", fg.Name),
		fmt.Sprintf("Entity:        %s", fg.EntityName),
		fmt.Sprintf("Description:   %s", fg.Description),
		fmt.Sprintf("Revision:      %s", revision),
		fmt.Sprintf("DataTable:     %s", dataTable),
		fmt.Sprintf("CreateTime:    %s", fg.CreateTime.Format(time.RFC3339)),
		fmt.Sprintf("ModifyTime:    %s", fg.ModifyTime.Format(time.RFC3339)),
	}, "\n")
}

func RichFeatureCsvHeader() string {
	return strings.Join([]string{"Name", "Group", "Entity", "Category", "DBValueType", "ValueType", "Description", "Revision", "DataTable", "CreateTime", "ModifyTime"}, ",")
}

func (r *RichFeature) ToCsvRecord() string {
	var revision, dataTable string
	if r.Revision == nil {
		revision = ""
	} else {
		revision = strconv.FormatInt(*r.Revision, 10)
	}
	if r.DataTable == nil {
		dataTable = ""
	} else {
		dataTable = *r.DataTable
	}

	return strings.Join([]string{r.Name, r.GroupName, r.EntityName, r.Category, r.DBValueType, r.ValueType, r.Description, revision, dataTable, r.CreateTime.Format(time.RFC3339), r.ModifyTime.Format(time.RFC3339)}, ",")
}

type RevisionRange struct {
	MinRevision int64  `db:"min_revision"`
	MaxRevision int64  `db:"max_revision"`
	DataTable   string `db:"data_table"`
}

type RawFeatureValueRecord struct {
	Record []interface{}
	Error  error
}

type EntityRow struct {
	EntityKey string `db:"entity_key"`
	UnixTime  int64  `db:"unix_time"`
}
