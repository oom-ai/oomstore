package database

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type FeatureConfig struct {
	Name           string    `db:"name"`
	Group          string    `db:"group"`
	Revision       string    `db:"revision"`
	Status         string    `db:"status"`
	Category       string    `db:"category"`
	ValueType      string    `db:"value_type"`
	Description    string    `db:"description"`
	RevisionsLimit int       `db:"revisions_limit"`
	CreateTime     time.Time `db:"create_time"`
	ModifyTime     time.Time `db:"modify_time"`
}

func (db *DB) ListFeatureConfig(ctx context.Context) ([]FeatureConfig, error) {
	query := `SELECT * FROM feature_config`
	features := make([]FeatureConfig, 0)
	if err := db.SelectContext(ctx, &features, query); err != nil {
		return nil, err
	}
	return features, nil
}

func (db *DB) ListFeatureConfigByGroup(ctx context.Context, group string) ([]FeatureConfig, error) {
	query := "SELECT * FROM feature_config AS fc WHERE fc.group = ?"
	features := make([]FeatureConfig, 0)
	if err := db.SelectContext(ctx, &features, query, group); err != nil {
		return nil, err
	}
	return features, nil
}

func (db *DB) GetFeatureConfig(ctx context.Context, groupName, featureName string) (*FeatureConfig, error) {
	var feature FeatureConfig
	query := `SELECT * FROM feature_config AS fc WHERE fc.group = ? AND fc.name = ?`
	if err := db.GetContext(ctx, &feature, query, groupName, featureName); err != nil {
		return nil, err
	}
	return &feature, nil
}

func (r *FeatureConfig) String() string {
	return strings.Join([]string{
		fmt.Sprintf("Name:           %s", r.Name),
		fmt.Sprintf("Group:          %s", r.Group),
		fmt.Sprintf("Revision:       %s", r.Revision),
		fmt.Sprintf("Status:         %s", r.Status),
		fmt.Sprintf("Category:       %s", r.Category),
		fmt.Sprintf("ValueType:      %s", r.ValueType),
		fmt.Sprintf("Description:    %s", r.Description),
		fmt.Sprintf("RevisionsLimit: %d", r.RevisionsLimit),
		fmt.Sprintf("CreateTime:     %s", r.CreateTime.Format(time.RFC3339)),
		fmt.Sprintf("ModifyTime:     %s", r.ModifyTime.Format(time.RFC3339)),
	}, "\n")
}

func (r *FeatureConfig) OneLineString() string {
	return strings.Join([]string{
		r.Name, r.Group, r.Revision, r.Status, r.Category, r.ValueType, r.Description,
		fmt.Sprintf("%d", r.RevisionsLimit),
		r.CreateTime.Format(time.RFC3339), r.ModifyTime.Format(time.RFC3339)},
		",")
}
