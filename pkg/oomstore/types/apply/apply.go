package apply

import "io"

type ApplyOpt struct {
	R io.Reader
}

type Feature struct {
	Kind        string `mapstructure:"kind"`
	Name        string `mapstructure:"name"`
	GroupID     int
	GroupName   string `mapstructure:"group-name"`
	DBValueType string `mapstructure:"db-type-value"`
	Description string `mapstructure:"description"`
}

type FeatureGroup struct {
	Kind        string    `mapstructure:"kind"`
	Group       string    `mapstructure:"group"`
	Name        string    `mapstructure:"name"`
	EntityName  string    `mapstructure:"entity_name"`
	Category    string    `mapstructure:"category"`
	Description string    `mapstructure:"description"`
	Features    []Feature `mapstructure:"features"`
}

type Entity struct {
	Kind        string `mapstructure:"kind"`
	Name        string `mapstructure:"name"`
	Length      int    `mapstructure:"length"`
	Description string `mapstructure:"description"`

	BatchFeatures []FeatureGroup `mapstructure:"batch-features"`
	StreamFeature []Feature      `mapstructure:"stream-features"`
}
