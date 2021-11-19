package apply

import "io"

type ApplyOpt struct {
	R io.Reader
}

type ApplyStage struct {
	NewEntities []Entity
	NewGroups   []Group
	NewFeatures []Feature
}

func NewApplyStage() *ApplyStage {
	return &ApplyStage{
		NewEntities: make([]Entity, 0),
		NewGroups:   make([]Group, 0),
		NewFeatures: make([]Feature, 0),
	}
}

type Feature struct {
	Kind        string `mapstructure:"kind"`
	Name        string `mapstructure:"name"`
	GroupID     int
	GroupName   string `mapstructure:"group-name"`
	DBValueType string `mapstructure:"db-value-type"`
	Description string `mapstructure:"description"`
}

type Group struct {
	Kind        string `mapstructure:"kind"`
	Group       string `mapstructure:"group"`
	Name        string `mapstructure:"name"`
	EntityID    int
	EntityName  string    `mapstructure:"entity-name"`
	Category    string    `mapstructure:"category"`
	Description string    `mapstructure:"description"`
	Features    []Feature `mapstructure:"features"`
}

type Entity struct {
	Kind        string `mapstructure:"kind"`
	Name        string `mapstructure:"name"`
	Length      int    `mapstructure:"length"`
	Description string `mapstructure:"description"`

	BatchFeatures  []Group   `mapstructure:"batch-features"`
	StreamFeatures []Feature `mapstructure:"stream-features"`
}
