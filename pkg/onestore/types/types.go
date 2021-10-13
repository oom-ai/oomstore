package types

import "time"

type Entity struct {
	Name        string `db:"name"`
	Description string `db:"description"`

	CreateTime time.Time `db:"create_time"`
	ModifyTime time.Time `db:"modify_time"`
}

type Feature struct {
	Name        string    `db:"name"`
	GroupName   string    `db:"group_name"`
	EntityName  string    `db:"entity_name"`
	ValueType   string    `db:"value_type"`
	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`
}

type Revision struct {
	Revision  string `db:"revision"`
	GroupName string `db:"group_name"`
	DataTable string `db:"data_table"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`
}
