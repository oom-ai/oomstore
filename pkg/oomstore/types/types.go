package types

import (
	"fmt"
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

type FeatureGroup struct {
	ID               int16  `db:"id"`
	Name             string `db:"name"`
	EntityName       string `db:"entity_name"`
	OnlineRevisionID *int32 `db:"online_revision_id"`
	Category         string `db:"category"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`
}

type RichFeatureGroup struct {
	FeatureGroup
	OnlineRevision   *int64  `db:"online_revision"`
	OfflineRevision  *int64  `db:"offline_revision"`
	OfflineDataTable *string `db:"offline_data_table"`
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

func (rfg *RichFeatureGroup) String() string {
	onlineRevision := "<NULL>"
	offlineRevision := "<NULL>"
	offlineDataTable := "<NULL>"

	if rfg.OnlineRevision != nil {
		onlineRevision = fmt.Sprint(*rfg.OnlineRevision)
	}
	if rfg.OfflineRevision != nil {
		offlineRevision = fmt.Sprint(*rfg.OfflineRevision)
	}
	if rfg.OfflineDataTable == nil {
		offlineDataTable = *rfg.OfflineDataTable
	}
	return strings.Join([]string{
		fmt.Sprintf("Name:                     %s", rfg.Name),
		fmt.Sprintf("Entity:                   %s", rfg.EntityName),
		fmt.Sprintf("Description:              %s", rfg.Description),
		fmt.Sprintf("Online Revision:          %s", onlineRevision),
		fmt.Sprintf("Offline Latest Revision:  %s", offlineRevision),
		fmt.Sprintf("Offline Latest DataTable: %s", offlineDataTable),
		fmt.Sprintf("CreateTime:               %s", rfg.CreateTime.Format(time.RFC3339)),
		fmt.Sprintf("ModifyTime:               %s", rfg.ModifyTime.Format(time.RFC3339)),
	}, "\n")
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

func (e *Entity) String() string {
	return strings.Join([]string{
		fmt.Sprintf("Name:          %s", e.Name),
		fmt.Sprintf("Length:        %d", e.Length),
		fmt.Sprintf("Description:   %s", e.Description),
		fmt.Sprintf("CreateTime:    %s", e.CreateTime.Format(time.RFC3339)),
		fmt.Sprintf("ModifyTime:    %s", e.ModifyTime.Format(time.RFC3339)),
	}, "\n")
}
