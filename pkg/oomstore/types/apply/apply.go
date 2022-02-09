package apply

import (
	"io"
	"time"

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
	Kind        string `mapstructure:"kind" yaml:"kind,omitempty" json:"kind,omitempty"`
	Name        string `mapstructure:"name" yaml:"name" json:"name"`
	GroupName   string `mapstructure:"group-name" yaml:"group-name,omitempty" json:"group-name,omitempty"`
	ValueType   string `mapstructure:"value-type" yaml:"value-type" json:"value-type"`
	Description string `mapstructure:"description" yaml:"description" json:"description"`
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

func BuildFeatureItems(features types.FeatureList) FeatureItems {
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
	Kind             string         `mapstructure:"kind" yaml:"kind,omitempty" json:"kind,omitempty"`
	Name             string         `mapstructure:"name" yaml:"name,omitempty" json:"name"`
	EntityName       string         `mapstructure:"entity-name" yaml:"entity-name,omitempty" json:"entity-name,omitempty"`
	Category         types.Category `mapstructure:"category" yaml:"category,omitempty" json:"category,omitempty"`
	SnapshotInterval time.Duration  `mapstructure:"snapshot-interval" yaml:"snapshot-interval,omitempty" json:"snapshot-interval"`
	Description      string         `mapstructure:"description" yaml:"description" json:"description"`
	Features         []Feature      `mapstructure:"features" yaml:"features,omitempty" json:"features,omitempty"`
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

func BuildGroupItems(groups types.GroupList, features types.FeatureList) *GroupItems {
	items := &GroupItems{
		Items: make([]Group, 0, groups.Len()),
	}

	for _, group := range groups {
		items.Items = append(items.Items, Group{
			Kind:             "Group",
			Name:             group.Name,
			EntityName:       group.Entity.Name,
			Category:         group.Category,
			SnapshotInterval: time.Duration(group.SnapshotInterval) * time.Second,
			Description:      group.Description,
			Features: BuildFeatureItems(features.Filter(func(f *types.Feature) bool {
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

type Entity struct {
	Kind        string `mapstructure:"kind" yaml:"kind" json:"kind,omitempty"`
	Name        string `mapstructure:"name" yaml:"name" json:"name"`
	Description string `mapstructure:"description" yaml:"description" json:"description"`

	Groups []Group `mapstructure:"groups" yaml:"groups,omitempty"`
}

type EntityItems struct {
	Items []Entity `yaml:"items"`
}

func BuildEntityItems(entities types.EntityList, groups *GroupItems) (items EntityItems) {
	for _, entity := range entities {
		items.Items = append(items.Items, Entity{
			Kind:        "Entity",
			Name:        entity.Name,
			Description: entity.Description,
			Groups: groups.Filter(func(g Group) bool {
				return g.EntityName == entity.Name
			}).Walk(func(g Group) Group {
				g.Kind = ""
				g.EntityName = ""
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
