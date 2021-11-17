package oomstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/oom-ai/oomstore/internal/database/metadata"
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

		group, err := s.metadata.GetGroupByName(ctx, feature.GroupName)
		if err != nil {
			return err
		}
		feature.GroupID = group.ID

		return s.metadata.WithTransaction(ctx, func(c context.Context, tx metadata.WriteStore) error {
			return s.applyFeature(ctx, tx, feature)
		})
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

func (s *OomStore) applyFeature(ctx context.Context, txStore metadata.WriteStore, newFeature apply.Feature) error {
	featureExist := true

	feature, err := s.metadata.GetFeatureByName(ctx, newFeature.Name)
	if err != nil {
		if err.Error() != fmt.Sprintf("feature '%s' not found", newFeature.Name) {
			return err
		}
		featureExist = false
	}

	valueType, err := s.offline.TypeTag(newFeature.DBValueType)
	if err != nil {
		return err
	}

	if !featureExist {
		_, err = txStore.CreateFeature(ctx, metadata.CreateFeatureOpt{
			FeatureName: newFeature.Name,
			GroupID:     newFeature.GroupID,
			DBValueType: newFeature.DBValueType,
			Description: newFeature.Description,
			ValueType:   valueType,
		})
		return err
	}

	if newFeature.Description != feature.Description {
		return txStore.UpdateFeature(ctx, metadata.UpdateFeatureOpt{
			FeatureID:      feature.ID,
			NewDescription: newFeature.Description,
		})
	}
	return nil
}
