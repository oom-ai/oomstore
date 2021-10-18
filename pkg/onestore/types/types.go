package types

import "time"

const (
	BatchFeatureCategory  = "batch"
	StreamFeatureCategory = "stream"
)

type Entity struct {
	Name   string `db:"name"`
	Length int    `db:"length"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`
}

type Feature struct {
	Name      string `db:"name"`
	GroupName string `db:"group_name"`
	ValueType string `db:"value_type"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`
}

type RichFeature struct {
	Feature
	EntityName string `db:"entity_name"`
	Category   string `db:"category"`
	Revision   int64  `db:"revision"`
	DataTable  string `db:"data_table"`
}

type Revision struct {
	Revision  int64  `db:"revision"`
	GroupName string `db:"group_name"`
	DataTable string `db:"data_table"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`
}

type FeatureGroup struct {
	Name       string `db:"name"`
	EntityName string `db:"entity_name"`
	Revision   int64  `db:"revision"`
	Category   string `db:"category"`
	DataTable  string `db:"data_table"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`
}

func (rf *RichFeature) ToFeature() *Feature {
	if rf == nil {
		return nil
	}
	return &rf.Feature
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

type FeatureDataSet map[string][]FeatureKV

func NewFeatureDataSet() FeatureDataSet {
	return make(map[string][]FeatureKV)
}
