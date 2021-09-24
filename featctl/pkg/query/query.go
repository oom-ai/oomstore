package query

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cast"

	"github.com/onestore-ai/onestore/featctl/pkg/database"
)

type Option struct {
	Group        string
	FeatureNames []string
	Entitykeys   []string
	Revision     string
	DBOption     database.Option
}

var firstPrint = true

func Run(ctx context.Context, opt *Option) {
	db, err := database.Open(&opt.DBOption)
	if err != nil {
		log.Fatalf("failed connecting feature store: %v", err)
	}
	defer db.Close()

	if err := queryFeatureAndPrintToStdout(ctx, db, opt); err != nil {
		log.Fatal(err)
	}
}

func queryFeatureAndPrintToStdout(ctx context.Context, db *database.DB, opt *Option) error {
	entityTableMapFeatures, err := getEntityTableMapFeatures(ctx, db, opt)
	if err != nil {
		return err
	}

	w := csv.NewWriter(os.Stdout)
	for entityTable, featureNames := range entityTableMapFeatures {
		if err := readOneTableToCsv(ctx, db, entityTable, opt.Entitykeys, featureNames, w); err != nil {
			return err
		}
	}
	return nil
}

func getEntityTableMapFeatures(ctx context.Context, db *database.DB, opt *Option) (map[string][]string, error) {
	mp := make(map[string][]string)

	if opt.Revision != "" {
		entityTable := opt.Group + "_" + opt.Revision
		mp[entityTable] = opt.FeatureNames
		return mp, nil
	}

	for _, featureName := range opt.FeatureNames {
		if entityTable, err := getEntityTable(ctx, db, opt.Group, featureName); err == nil && entityTable != "" {
			if v, ok := mp[entityTable]; ok {
				mp[entityTable] = append(v, featureName)
			} else {
				mp[entityTable] = []string{featureName}
			}
		} else {
			return nil, fmt.Errorf("cannot find entity table for group=%s, featureName=%s, err: %v", opt.Group, featureName, err)
		}
	}
	return mp, nil
}

func getEntityTable(ctx context.Context, db *database.DB, group, featureName string) (string, error) {
	var revision string
	err := db.QueryRowContext(ctx, `select fc.revision from feature_config as fc where fc.group = ? and fc.name = ?`, group, featureName).Scan(&revision)
	switch {
	case err == sql.ErrNoRows:
		return "", nil
	case err != nil:
		return "", err
	default:
		return group + "_" + revision, nil
	}
}

func readOneTableToCsv(ctx context.Context, db *database.DB, tableName string,
	entityKeys []string, featureNames []string, w *csv.Writer) error {
	sql := fmt.Sprintf("select entity_key, %s from %s", strings.Join(featureNames, ", "), tableName)
	if len(entityKeys) > 0 {
		sql += fmt.Sprintf(" where entity_key in (%s)", strings.Join(entityKeys, ", "))
	}

	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("failed connecting feature store: %v", err)
	}
	defer rows.Close()

	return resolveDataFromRows(rows, w)
}

func resolveDataFromRows(rows *sql.Rows, w *csv.Writer) error {
	if rows == nil {
		return fmt.Errorf("rows can't be nil")
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	length := len(columns)

	record := make([]string, length, length)
	if firstPrint {
		for i, column := range columns {
			record[i] = column
		}
		w.Write(record)
		firstPrint = false
	}
	//unnecessary to put below into rows.Next loop,reduce allocating
	values := make([]interface{}, length)
	for i := 0; i < length; i++ {
		values[i] = new(interface{})
	}

	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			return err
		}

		for i := 0; i < len(columns); i++ {
			value := *(values[i].(*interface{}))
			record[i] = cast.ToString(value)
		}

		w.Write(record)
	}
	w.Flush()
	return nil
}
