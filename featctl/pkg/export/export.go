package export

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/joho/sqltocsv"
	"github.com/onestore-ai/onestore/pkg/database"
)

type Option struct {
	Group      string
	Features   []string
	OutputFile string
	DBOption   database.Option
}

func downloadFeatures(ctx context.Context, db *database.DB, opt *Option, tableName string) error {
	dbo := opt.DBOption
	fields := strings.Join(opt.Features, ",")

	if len(opt.Features) == 0 {
		// download all fields by default
		fields = "*"
	} else if !containsField(opt.Features, "*") && !containsField(opt.Features, "entity_key") {
		// make sure the field `entity_key` is included
		fields = "entity_key," + fields
	}

	fullTableName := fmt.Sprintf("%s.%s", dbo.DbName, tableName)
	query := fmt.Sprintf("select %s from %s", fields, fullTableName)

	return dumpCSV(ctx, db, opt.OutputFile, query)
}

func dumpCSV(ctx context.Context, db *database.DB, file string, query string, args ...interface{}) error {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return sqltocsv.WriteFile(file, rows)
}

func Export(ctx context.Context, option *Option) {
	log.Println("connecting feature store ...")
	db, err := database.Open(&option.DBOption)
	if err != nil {
		log.Fatalf("failed connecting feature store: %v", err)
	}
	defer db.Close()

	log.Println("retrieving source table ...")
	sourceTableName, err := database.GetLatestEntityTable(ctx, db, option.Group)
	if err != nil {
		log.Fatalf("failed retrieving source table: %v", err)
	}

	log.Println("downloading features ...")
	if err = downloadFeatures(ctx, db, option, sourceTableName); err != nil {
		log.Fatalf("failed downloading features: %v", err)
	}

	log.Println("succeeded.")
}

func containsField(fields []string, target string) bool {
	target = strings.TrimSpace(target)
	for _, s := range fields {
		if strings.TrimSpace(s) == target {
			return true
		}
	}
	return false
}
