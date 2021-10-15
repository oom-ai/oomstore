package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func (db *DB) LoadLocalFile(ctx context.Context, filePath, tableName, separator, delimiter string, header []string) error {
	loadData := fmt.Sprintf(
		"LOAD DATA LOCAL INFILE '%s' INTO TABLE %s FIELDS TERMINATED BY '%s' ENCLOSED BY '%s' LINES TERMINATED BY '\n' IGNORE 1 LINES (%s)",
		filePath,
		tableName,
		separator,
		delimiter,
		strings.Join(header, ","))
	mysql.RegisterLocalFile(filePath)
	_, err := db.ExecContext(ctx, loadData)
	return err
}
