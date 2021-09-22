package describe

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/onestore-ai/onestore/featctl/pkg/database"
)

type Option struct {
	Name     string
	Group    string
	DBOption database.Option
}

type FeatureConfigRecord struct {
	Name           string
	Group          string
	Revision       string
	Status         string
	Category       string
	ValueType      string
	Description    string
	RevisionsLimit int
	CreateTime     time.Time
	ModifyTime     time.Time
}

func Run(ctx context.Context, option *Option) {
	db, err := database.Open(&option.DBOption)
	if err != nil {
		log.Fatalf("failed connecting feature store: %v", err)
	}
	defer db.Close()

	record, err := describe(ctx, db, option)
	if err != nil {
		log.Fatalf("failed querying feature config: %v", err)
	}
	if record == nil {
		fmt.Fprintln(os.Stderr, "Feature not found.")
		os.Exit(4)
	}

	fmt.Println(record.String())
}

func describe(ctx context.Context, db *database.DB, option *Option) (*FeatureConfigRecord, error) {
	var record FeatureConfigRecord
	err := db.QueryRowContext(ctx,
		`select fc.name,
				fc.group,
				fc.revision,
				fc.status,
				fc.category,
				fc.value_type,
				fc.description,
				fc.revisions_limit,
				fc.create_time,
				fc.modify_time
			from feature_config as fc
			where fc.group = ? and fc.name = ?`,
		option.Group,
		option.Name,
	).Scan(
		&record.Name,
		&record.Group,
		&record.Revision,
		&record.Status,
		&record.Category,
		&record.ValueType,
		&record.Description,
		&record.RevisionsLimit,
		&record.CreateTime,
		&record.ModifyTime,
	)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return &record, nil
	}
}

func (r *FeatureConfigRecord) String() string {
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
