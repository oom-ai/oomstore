package export

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/onestore-ai/onestore/pkg/database"
)

type Option struct {
	Group      string
	Features   []string
	OutputFile string
	DBOption   database.Option
}

func getSourceTableName(ctx context.Context, db *database.DB, group string) (string, error) {
	var source string
	err := db.QueryRowContext(ctx,
		"select source from feature_revision"+
			" where `group` = ?"+
			" order by revision desc"+
			" limit 1", group).
		Scan(&source)
	if err != nil {
		return "", err
	}
	return source, nil
}

func downloadFeatures(ctx context.Context, opt *Option, tableName string) error {
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
	sql := fmt.Sprintf("select %s from %s", fields, fullTableName)
	tmpdir := opt.OutputFile + ".tmp"

	// download data using tidb-dumpling
	// https://docs.pingcap.com/zh/tidb/v4.0/dumpling-overview
	cmd := exec.CommandContext(ctx,
		"dumpling",
		"-h", dbo.Host,
		"-P", dbo.Port,
		"-u", dbo.User,
		"-p", dbo.Pass,
		"-T", fullTableName,
		"-o", tmpdir,
		"--filetype", "csv",
		"--escape-backslash=false",
		"--consistency", "none",
		"--tidb-mem-quota-query", strconv.Itoa(512*1024*1024), // 512M
		"--params", "tidb_distsql_scan_concurrency=1", // lowest scan concurrency
		"--sql", sql,
	)

	cmdOutput, err := cmd.CombinedOutput()
	fmt.Println(string(cmdOutput))
	if err != nil {
		return err
	}

	// move csv file to the specified location
	cmd = exec.CommandContext(ctx,
		"sh",
		"-c",
		fmt.Sprintf("mv %s/*.csv %s", tmpdir, opt.OutputFile),
	)
	cmdOutput, err = cmd.CombinedOutput()
	fmt.Println(string(cmdOutput))
	if err != nil {
		return err
	}

	// remove temporary directory
	return os.RemoveAll(tmpdir)
}

func Export(ctx context.Context, option *Option) {
	log.Println("connecting feature store ...")
	db, err := database.Open(&option.DBOption)
	if err != nil {
		log.Fatalf("failed connecting feature store: %v", err)
	}
	defer db.Close()

	log.Println("retrieving source table ...")
	sourceTableName, err := getSourceTableName(ctx, db, option.Group)
	if err != nil {
		log.Fatalf("failed retrieving source table: %v", err)
	}

	log.Println("downloading features ...")
	if err = downloadFeatures(ctx, option, sourceTableName); err != nil {
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
