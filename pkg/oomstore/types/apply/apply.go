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
	ValueType   string `mapstructure:"value-type" yaml:"value-type"`
	Description string `mapstructure:"description" yaml:"description"`
}

func (f *Feature) Validate() error {
	if f.Name == "" {
		return errdefs.InvalidAttribute(fmt.Errorf("the name of feature should not be empty"))
	}
	if f.ValueType == "" {
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
			ValueType:   f.ValueType.String(),
			Description: f.Description,
		})
	}
	return items
}

type Group struct {
	Kind        string    `mapstructure:"kind" yaml:"kind,omitempty"`
	Group       string    `mapstructure:"group" yaml:"group,omitempty"`
	Name        string    `mapstructure:"name" yaml:"name,omitempty"`
	EntityName  string    `mapstructure:"entity-name" yaml:"entity-name,omitempty"`
	Category    string    `mapstructure:"category" yaml:"category,omitempty"`
	Description string    `mapstructure:"description" yaml:"description"`
	Features    []Feature `mapstructure:"features" yaml:"features,omitempty"`
}

type GroupItems struct {
	Items []Group `mapstructure:"items" yaml:"items"`
}

func (g *GroupItems) Filter(filter func(Group) bool) (rs GroupItems) {
	for _, item := range g.Items {
		if filter(item) {
			rs.Items = append(rs.Items, item)
		}
	}
	return
}

func (g GroupItems) Walk(walkFunc func(Group) Group) (rs GroupItems) {
	for _, i := range g.Items {
		rs.Items = append(rs.Items, walkFunc(i))
	}
	return
}

func FromGroupList(groups types.GroupList, features types.FeatureList) *GroupItems {
	items := &GroupItems{
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
	Kind        string `mapstructure:"kind" yaml:"kind"`
	Name        string `mapstructure:"name" yaml:"name"`
	Length      int    `mapstructure:"length" yaml:"length"`
	Description string `mapstructure:"description" yaml:"description"`

	BatchFeatures  []Group   `mapstructure:"batch-features" yaml:"batch-features,omitempty"`
	StreamFeatures []Feature `mapstructure:"stream-features" yaml:"stream-features,omitempty"`
}

func (e *Entity) Validate() error {
	if e.Name == "" {
		return errdefs.InvalidAttribute(fmt.Errorf("the name of entity should not be empty"))
	}
	return nil
}

type EntityItems struct {
	Items []Entity `yaml:"items"`
}

func FromEntityList(entities types.EntityList, groups *GroupItems) (items EntityItems) {
	for _, entity := range entities {
		items.Items = append(items.Items, Entity{
			Kind:        "Entity",
			Name:        entity.Name,
			Length:      entity.Length,
			Description: entity.Description,
			BatchFeatures: groups.Filter(func(g Group) bool {
				return g.EntityName == entity.Name
			}).Walk(func(g Group) Group {
				g.Kind = ""
				g.EntityName = ""
				g.Category = ""
				g.Group = g.Name
				g.Name = ""
				return g
			}).Items,
		})
	}
	return
}

type Kind struct {
	Kind string `mapstructure:"kind"`
}

type Items struct {
	Items []Kind `mapstructure:"items"`
}

func (i *Items) Kind() string {
	if len(i.Items) == 0 {
		return ""
	}
	return i.Items[0].Kind
}
