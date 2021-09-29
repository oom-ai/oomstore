package _import

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/onestore-ai/onestore/pkg/database"
)

type Option struct {
	Group          string
	Revision       string
	SchemaTemplate string
	InputOption    InputOption
	DBOption       database.Option
	Description    string
}

type InputOption struct {
	FilePath  string
	NoHeader  bool
	Separator string
	Delimiter string
}

func Run(ctx context.Context, opt *Option) {
	db, err := database.Open(&opt.DBOption)
	if err != nil {
		log.Fatalf("failed connecting feature store: %v", err)
	}
	defer db.Close()

	if err := validateOptions(ctx, db, opt); err != nil {
		log.Fatalf("failed validating options: %v\n", err)
	}

	log.Println("importing features ...")
	if err := importFeatures(ctx, db, opt); err != nil {
		log.Fatalf("failed importing features: %v\n", err)
	}

	log.Println("registering revision ...")
	if err := database.RegisterRevision(ctx, db,
		opt.Group,
		opt.Revision,
		genTableName(opt),
		opt.Description,
	); err != nil {
		log.Fatalf("failed registering revision: %v\n", err)
	}

	log.Println("succeeded.")
}

func genTableName(opt *Option) string {
	return fmt.Sprintf("%s_%s", opt.Group, opt.Revision)
}

func importFeatures(ctx context.Context, db *database.DB, opt *Option) error {
	// create table
	schema, err := genSchema(opt)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, schema)
	if err != nil {
		return err
	}

	// import data
	table := genTableName(opt)
	iopt := opt.InputOption
	stmt := fmt.Sprintf(`LOAD DATA LOCAL INFILE '%s' INTO TABLE %s FIELDS TERMINATED BY '%s' ENCLOSED BY '%s' LINES TERMINATED BY '\n'`,
		iopt.FilePath, table, iopt.Separator, iopt.Delimiter)
	if !iopt.NoHeader {
		stmt += " IGNORE 1 LINES"
	}

	mysql.RegisterLocalFile(iopt.FilePath)
	_, err = db.ExecContext(ctx, stmt)
	return err
}

func genSchema(opt *Option) (string, error) {
	bytes, err := os.ReadFile(opt.SchemaTemplate)
	if err != nil {
		return "", err
	}

	schema := string(bytes)
	if !strings.Contains(schema, "{{TABLE_NAME}}") {
		return "", fmt.Errorf("'{{TABLE_NAME}}' not found in schema template", opt.SchemaTemplate)
	}
	schema = strings.ReplaceAll(schema, "{{TABLE_NAME}}", genTableName(opt))
	return schema, nil
}

func validateOptions(ctx context.Context, db *database.DB, option *Option) error {
	table := genTableName(option)
	exists, err := db.TableExists(ctx, genTableName(option))
	if err != nil {
		return fmt.Errorf("failed querying database => %v", err)
	}
	if exists {
		return fmt.Errorf("invalid group or revision => table '%s' already exists", table)
	}
	return nil
}

func registerRevision(ctx context.Context, db *database.DB, option *Option) error {
	_, err := db.ExecContext(ctx,
		"insert into feature_revision(`group`, revision, source, description) values(?, ?, ?, ?)",
		option.Group,
		option.Revision,
		genTableName(option),
		option.Description)
	return err
}
