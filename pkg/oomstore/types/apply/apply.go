package apply

import (
	"fmt"
	"io"

	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

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
	Kind        string `mapstructure:"kind" yaml:"kind,omitempty"`
	Name        string `mapstructure:"name" yaml:"name"`
	GroupName   string `mapstructure:"group-name" yaml:"group-name,omitempty"`
	DBValueType string `mapstructure:"db-value-type" yaml:"db-value-type"`
	Description string `mapstructure:"description" yaml:"description"`
}

func (f *Feature) Validate() error {
	if f.Name == "" {
		return errdefs.InvalidAttribute(fmt.Errorf("the name of feature should not be empty"))
	}
	if f.DBValueType == "" {
		return errdefs.InvalidAttribute(fmt.Errorf("the db value type of feature should not be empty"))
	}
	return nil
}

type FeatureItems struct {
	Items []Feature `mapstructure:"items" yaml:"items"`
}

func (f FeatureItems) Walk(walkFunc func(Feature) Feature) (rs FeatureItems) {
	for _, i := range f.Items {
		rs.Items = append(rs.Items, walkFunc(i))
	}
	return
}

func FromFeatureList(features types.FeatureList) FeatureItems {
	items := FeatureItems{
		Items: make([]Feature, 0, features.Len()),
	}

	for _, f := range features {
		items.Items = append(items.Items, Feature{
			Kind:        "Feature",
			Name:        f.Name,
			GroupName:   f.Group.Name,
			DBValueType: f.DBValueType,
			Description: f.Description,
		})
	}
	return items
}

type Group struct {
	Kind        string    `mapstructure:"kind" yaml:"kind"`
	Group       string    `mapstructure:"group" yaml:"group,omitempty"`
	Name        string    `mapstructure:"name" yaml:"name,omitempty"`
	EntityName  string    `mapstructure:"entity-name" yaml:"entity-name"`
	Category    string    `mapstructure:"category" yaml:"category"`
	Description string    `mapstructure:"description" yaml:"description"`
	Features    []Feature `mapstructure:"features" yaml:"features"`
}

type GroupItems struct {
	Items []Group `mapstructure:"item" yaml:"items"`
}

func FromGroupList(groups types.GroupList, features types.FeatureList) GroupItems {
	items := GroupItems{
		Items: make([]Group, 0, groups.Len()),
	}

	for _, group := range groups {
		items.Items = append(items.Items, Group{
			Kind:        "Group",
			Name:        group.Name,
			EntityName:  group.Entity.Name,
			Category:    group.Category,
			Description: group.Description,
			Features: FromFeatureList(features.Filter(func(f *types.Feature) bool {
				return f.Group.Name == group.Name
			})).Walk(func(f Feature) Feature {
				f.Kind = ""
				f.GroupName = ""
				return f
			}).Items,
		})
	}

	return items
}

func (g *Group) Validate() error {
	if g.Name == "" {
		return errdefs.InvalidAttribute(fmt.Errorf("the name of group should not be empty"))
	}
	return nil
}

type Entity struct {
	Kind        string `mapstructure:"kind"`
	Name        string `mapstructure:"name"`
	Length      int    `mapstructure:"length"`
	Description string `mapstructure:"description"`

	BatchFeatures  []Group   `mapstructure:"batch-features"`
	StreamFeatures []Feature `mapstructure:"stream-features"`
}

func (e *Entity) Validate() error {
	if e.Name == "" {
		return errdefs.InvalidAttribute(fmt.Errorf("the name of entity should not be empty"))
	}
	return nil
}
