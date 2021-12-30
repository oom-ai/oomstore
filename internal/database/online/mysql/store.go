package mysql

import (
	"context"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const BackendType = types.BackendMySQL

var _ online.Store = &DB{}

type DB struct {
	*sqlx.DB
}

func Open(opt *types.MySQLOpt) (*DB, error) {
	db, err := dbutil.OpenMysqlDB(opt.Host, opt.Port, opt.User, opt.Password, opt.Database)
	return &DB{db}, err
}

func (db *DB) Ping(ctx context.Context) error {
	return db.PingContext(ctx)
}

func (db *DB) Get(ctx context.Context, opt online.GetOpt) (dbutil.RowMap, error) {
	return sqlutil.Get(ctx, db.DB, opt, BackendType)
}

func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	return sqlutil.MultiGet(ctx, db.DB, opt, BackendType)
}

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	return sqlutil.Import(ctx, db.DB, opt, BackendType)
}

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	return sqlutil.Purge(ctx, db.DB, revisionID, BackendType)
}

// TODO: refactor with text/template
func (db *DB) Push(ctx context.Context, opt online.PushOpt) error {
	tableName := sqlutil.OnlineStreamTableName(opt.GroupID)

	insertColumns := append([]string{opt.Entity.Name}, opt.FeatureNames...)
	insertColumnPlaceholders := make([]string, 0, len(insertColumns))
	for i := 0; i < len(insertColumnPlaceholders); i++ {
		insertColumnPlaceholders = append(insertColumnPlaceholders, "?")
	}

	insertValues := append([]interface{}{opt.EntityKey}, opt.FeatureValues...)
	insertValuePlaceholders := make([]string, 0, len(insertValues))
	for i := 0; i < len(insertValues); i++ {
		insertValuePlaceholders = append(insertValuePlaceholders, "?")
	}

	updateValues := opt.FeatureValues
	updatePlaceholders := make([]string, 0, len(opt.FeatureNames))
	for _, name := range opt.FeatureNames {
		updatePlaceholders = append(updatePlaceholders, fmt.Sprintf("%s=?", name))
	}

	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES(%s) ON DUPLICATE KEY UPDATE %s`,
		tableName,
		strings.Join(insertColumns, ","),
		strings.Join(insertValuePlaceholders, ","),
		strings.Join(updatePlaceholders, ","),
	)
	_, err := db.ExecContext(ctx, query, append(insertValues, updateValues...)...)
	return err
}

func (db *DB) PrepareStreamTable(ctx context.Context, opt online.PrepareStreamTableOpt) error {
	return sqlutil.SqlxPrapareStreamTable(ctx, db.DB, opt, types.BackendMySQL)
}
