package types

import (
	"fmt"
	"strings"
	"time"
)

type Feature struct {
	ID          int16     `db:"id"`
	Name        string    `db:"name"`
	ValueType   string    `db:"value_type"`
	DBValueType string    `db:"db_value_type"`
	Category    string    `db:"category"`
	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`

	GroupId   int16  `db:"group_id"`
	GroupName string `db:"group_name"`

	EntityId   int16  `db:"entity_id"`
	EntityName string `db:"entity_name"`

	OnlineRevisionID *int32  `db:"online_revision_id"`
	OnlineRevision   *int64  `db:"online_revision"`
	OfflineRevision  *int64  `db:"offline_revision"`
	OfflineDataTable *string `db:"offline_data_table"`
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

func (l *FeatureList) Filter(filter func(f *Feature) bool) (rs FeatureList) {
	for _, f := range *l {
		if filter(f) {
			rs = append(rs, f)
		}
	}
	return
}

func (rf *Feature) String() string {
	onlineRevision := "<NULL>"
	offlineRevision := "<NULL>"
	offlineDataTable := "<NULL>"

	if rf.OnlineRevision != nil {
		onlineRevision = fmt.Sprint(*rf.OnlineRevision)
	}
	if rf.OfflineRevision != nil {
		offlineRevision = fmt.Sprint(*rf.OfflineRevision)
	}
	if rf.OfflineDataTable == nil {
		offlineDataTable = *rf.OfflineDataTable
	}

	return strings.Join([]string{
		fmt.Sprintf("Name:                     %s", rf.Name),
		fmt.Sprintf("Group:                    %s", rf.GroupName),
		fmt.Sprintf("Entity:                   %s", rf.EntityName),
		fmt.Sprintf("Category:                 %s", rf.Category),
		fmt.Sprintf("DBValueType:              %s", rf.DBValueType),
		fmt.Sprintf("ValueType:                %s", rf.ValueType),
		fmt.Sprintf("Description:              %s", rf.Description),
		fmt.Sprintf("Online Revision:          %s", onlineRevision),
		fmt.Sprintf("Offline Latest Revision:  %s", offlineRevision),
		fmt.Sprintf("Offline Latest DataTable:  %s", offlineDataTable),
		fmt.Sprintf("CreateTime:               %s", rf.CreateTime.Format(time.RFC3339)),
		fmt.Sprintf("ModifyTime:               %s", rf.ModifyTime.Format(time.RFC3339)),
	}, "\n")
}

func FeatureCsvHeader() string {
	return strings.Join([]string{"Name", "Group", "Entity", "Category", "DBValueType", "ValueType", "Description", "OnlineRevision", "OfflineLatestRevision", "OfflineLatestDataTable", "CreateTime", "ModifyTime"}, ",")
}

func (rf *Feature) ToCsvRecord() string {
	onlineRevision := "<NULL>"
	offlineRevision := "<NULL>"
	offlineDataTable := "<NULL>"

	if rf.OnlineRevision != nil {
		onlineRevision = fmt.Sprint(*rf.OnlineRevision)
	}
	if rf.OfflineRevision != nil {
		offlineRevision = fmt.Sprint(*rf.OfflineRevision)
	}
	if rf.OfflineDataTable != nil {
		offlineDataTable = *rf.OfflineDataTable
	}

	return strings.Join([]string{rf.Name, rf.GroupName, rf.EntityName, rf.Category, rf.DBValueType, rf.ValueType, rf.Description, onlineRevision, offlineRevision, offlineDataTable, rf.CreateTime.Format(time.RFC3339), rf.ModifyTime.Format(time.RFC3339)}, ",")
}
