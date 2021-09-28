package _import

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

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
	FilePath          string
	HasHeader         bool
	Separator         string
	Delimiter         string
	NullLiteral       string
	BackslashEscape   bool
	TrimLastSeparator bool
}

const (
	lightningLog = "lightning.log"
	lightningCfg = "lightning.toml"
)

func Import(ctx context.Context, option *Option) {
	db, err := database.Open(&option.DBOption)
	if err != nil {
		log.Fatalf("failed connecting feature store: %v", err)
	}
	defer db.Close()

	log.Println("validating options ...")
	if err := validateOptions(ctx, db, option); err != nil {
		log.Fatalf("failed validating options: %v\n", err)
	}

	log.Println("preparing files for lightning ...")
	if err = prepareLightningFiles(option); err != nil {
		log.Fatalf("failed preparing lightning files: %v\n", err)
	}

	log.Println("importing using lightning ...")
	if err = lightningImport(ctx); err != nil {
		log.Fatalf("failed running lightning: %v\n", err)
	}

	log.Println("registering revision ...")
	if err = registerRevision(ctx, db, option); err != nil {
		log.Fatalf("failed registering revision: %v\n", err)
	}

	log.Println("succeeded.")
}

func genTableName(option *Option) string {
	return fmt.Sprintf("%s_%s", option.Group, option.Revision)
}

func genSchema(option *Option) (string, error) {
	bytes, err := os.ReadFile(option.SchemaTemplate)
	if err != nil {
		return "", err
	}

	schema := string(bytes)
	if !strings.Contains(schema, "{{TABLE_NAME}}") {
		return "", fmt.Errorf("'{{TABLE_NAME}}' not found in schema template", option.SchemaTemplate)
	}
	schema = strings.ReplaceAll(schema, "{{TABLE_NAME}}", genTableName(option))
	return schema, nil
}

func genLightningCfg(option *Option, dataSourceDir string) string {
	cfg := lightningCfgTemplate
	cfg = strings.ReplaceAll(cfg, "{{DATA_SOURCE_DIR}}", dataSourceDir)
	cfg = strings.ReplaceAll(cfg, "{{LOG_PATH}}", lightningLog)

	ifo := option.InputOption
	cfg = strings.ReplaceAll(cfg, "{{SEPARATOR}}", ifo.Separator)
	cfg = strings.ReplaceAll(cfg, "{{DELIMITER}}", ifo.Delimiter)
	cfg = strings.ReplaceAll(cfg, "{{NULL_LITERAL}}", ifo.NullLiteral)
	cfg = strings.ReplaceAll(cfg, "{{HAS_HEADER}}", strconv.FormatBool(ifo.HasHeader))
	cfg = strings.ReplaceAll(cfg, "{{BACK_SLASH_ESCAPE}}", strconv.FormatBool(ifo.BackslashEscape))
	cfg = strings.ReplaceAll(cfg, "{{TRIM_LAST_SEPARATOR}}", strconv.FormatBool(ifo.TrimLastSeparator))

	dbo := option.DBOption
	cfg = strings.ReplaceAll(cfg, "{{HOST}}", dbo.Host)
	cfg = strings.ReplaceAll(cfg, "{{PORT}}", dbo.Port)
	cfg = strings.ReplaceAll(cfg, "{{USER}}", dbo.User)
	cfg = strings.ReplaceAll(cfg, "{{PASS}}", dbo.Pass)
	return cfg
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

func prepareLightningFiles(option *Option) error {
	pathLightningData := fmt.Sprintf("%s.tmp/", option.InputOption.FilePath)

	// remove log and data dir if exists
	for _, dir := range []string{lightningLog, pathLightningData} {
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}
	// create data dir
	if err := os.Mkdir(pathLightningData, 0755); err != nil {
		return err
	}

	dbName := option.DBOption.DbName
	tableName := genTableName(option)

	pathData := filepath.Join(pathLightningData, fmt.Sprintf("%s.%s.csv", dbName, tableName))
	pathTableSchema := filepath.Join(pathLightningData, fmt.Sprintf("%s.%s-schema.sql", dbName, tableName))
	pathDBSchema := filepath.Join(pathLightningData, fmt.Sprintf("%s-schema-create.sql", dbName))

	// database create schema - we don't need it but lightning requires it
	if _, err := os.Create(pathDBSchema); err != nil {
		return err
	}

	// table create schema
	schemaContent, err := genSchema(option)
	if err != nil {
		return err
	}
	if err := os.WriteFile(pathTableSchema, []byte(schemaContent), 0644); err != nil {
		return err
	}

	// data file
	if err := os.Symlink(filepath.Join("..", option.InputOption.FilePath), pathData); err != nil {
		return err
	}

	// lightning config
	cfgContent := genLightningCfg(option, pathLightningData)
	if err := os.WriteFile(lightningCfg, []byte(cfgContent), 0644); err != nil {
		return err
	}

	return nil
}

func lightningImport(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "tidb-lightning", "-config", lightningCfg)

	cmdOutput, err := cmd.CombinedOutput()
	fmt.Println(string(cmdOutput))
	if err != nil {
		return err
	}

	// output logs
	cmd = exec.CommandContext(ctx, "cat", lightningLog)
	cmd.Stdout = os.Stdout
	return cmd.Run()
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

const lightningCfgTemplate = `
[lightning]
level = "info"
file = "{{LOG_PATH}}"

[tikv-importer]
backend = "tidb"
on-duplicate = "replace"

[mydumper]
data-source-dir = "{{DATA_SOURCE_DIR}}"

[mydumper.csv]
separator = "{{SEPARATOR}}"
delimiter = '{{DELIMITER}}'
null = '{{NULL_LITERAL}}'
header = {{HAS_HEADER}}
backslash-escape = {{BACK_SLASH_ESCAPE}}
trim-last-separator = {{TRIM_LAST_SEPARATOR}}
not-null = false

[tidb]
host = "{{HOST}}"
port = {{PORT}}
user = "{{USER}}"
password = "{{PASS}}"

[checkpoint]
enable = false

[post-restore]
checksum = false
`
