package oomstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/oom-ai/oomstore/pkg/oomstore/types/apply"
	"gopkg.in/yaml.v3"
)

func (s *OomStore) Apply(ctx context.Context, opt apply.ApplyOpt) error {
	data := make(map[string]interface{})
	if err := yaml.NewDecoder(opt.R).Decode(data); err != nil {
		return err
	}

	var kind string
	if k, ok := data["kind"]; ok {
		kind = k.(string)
	} else {
		return fmt.Errorf("Invalid data missing variable kind")
	}

	switch strings.ToLower(kind) {
	case "entity":
		var entity apply.Entity
		if err := mapstructure.Decode(data, &entity); err != nil {
			return err
		}
		return s.applyEntity(ctx, entity)
	case "feature":
		var feature apply.Feature
		if err := mapstructure.Decode(data, &feature); err != nil {
			return err
		}
		return s.applyFeature(ctx, feature)
	case "group":
		var group apply.FeatureGroup
		if err := mapstructure.Decode(data, &group); err != nil {
			return err
		}
		return s.applyGroup(ctx, group)
	default:
		return fmt.Errorf("invalid kind %s", kind)
	}
}

func (s *OomStore) applyEntity(ctx context.Context, entity apply.Entity) error {
	fmt.Printf("apply entity!")
	return nil
}

func (s *OomStore) applyGroup(ctx context.Context, group apply.FeatureGroup) error {
	fmt.Printf("apply group!")
	return nil
}

func (s *OomStore) applyFeature(ctx context.Context, feature apply.Feature) error {
	fmt.Printf("apply feature!")
	return nil
}
