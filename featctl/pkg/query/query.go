package query

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/onestore-ai/onestore/pkg/database"
)

type Option struct {
	Group        string
	FeatureNames []string
	Entitykeys   []string
	Revision     string
	DBOption     database.Option
}

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

	isFirstPrint := true
	w := csv.NewWriter(os.Stdout)
	for entityTable, featureNames := range entityTableMapFeatures {
		if err := readOneTableToCsv(ctx, db, entityTable, opt.Entitykeys, featureNames, w, isFirstPrint); err != nil {
			return err
		}
		isFirstPrint = false
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
		if entityTable, err := database.GetEntityTable(ctx, db, opt.Group, featureName); err == nil && entityTable != "" {
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

func readOneTableToCsv(ctx context.Context, db *database.DB, tableName string,
	entityKeys []string, featureNames []string, w *csv.Writer, isFirstPrint bool) error {
	rows, err := database.ReadEntityTable(ctx, db, tableName, entityKeys, featureNames)
	if err != nil {
		return err
	}
	defer rows.Close()

	return database.ReadRowsToCsvFile(rows, w, isFirstPrint)
}
