package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
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

func ListFeatureConfigByGroup(db *sqlx.DB, group string) ([]FeatureConfig, error) {
	query := fmt.Sprintf(`SELECT * FROM feature_config AS fc WHERE fc.group = '%s'`, group)
	features := make([]FeatureConfig, 0)
	if err := db.Select(&features, query); err != nil {
		return nil, err
	}
	return features, nil
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
