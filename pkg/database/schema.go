package database

import (
	"context"
	"fmt"
)

func CreateFeatureRevisionTable(ctx context.Context, db *DB) error {
	exist, err := db.TableExists(ctx, "feature_revision")
	if err != nil || exist {
		return err
	}

	schema :=
		"CREATE TABLE feature_revision (" +
			"  `group`       VARCHAR(32)  NOT     NULL," +
			"  `revision`    VARCHAR(64)  NOT     NULL," +
			"  `source`      VARCHAR(64)  NOT     NULL," +
			"  `description` VARCHAR(128) DEFAULT NULL," +
			"  `create_time` TIMESTAMP    NOT     NULL DEFAULT CURRENT_TIMESTAMP," +
			"  `modify_time` TIMESTAMP    NOT     NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
			"  PRIMARY KEY pk_feature_revision(`group`, `revision`)" +
			")"
	_, err = db.ExecContext(ctx, schema)
	return err
}

func CreateFeatureConfigTable(ctx context.Context, db *DB) error {
	exist, err := db.TableExists(ctx, "feature_config")
	if err != nil || exist {
		return err
	}

	schema :=
		"CREATE TABLE feature_config (" +
			"  `group`			 VARCHAR(128) NOT NULL," +
			"  `name`            VARCHAR(64)  NOT NULL," +
			"  `revision`        VARCHAR(64)  NOT NULL," +
			"  `status`          VARCHAR(32)  NOT NULL," +
			"  `category`        VARCHAR(16)  NOT NULL," +
			"  `value_type`      VARCHAR(32)  NOT NULL," +
			"  `description`     VARCHAR(64)  NOT NULL," +
			"  `revisions_limit` INT          NOT NULL," +
			"  `create_time`     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP," +
			"  `modify_time`     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
			"  PRIMARY KEY pk_feature_config(`name`, `group`)" +
			")"
	_, err = db.ExecContext(ctx, schema)
	return err
}

func CreateDatabase(ctx context.Context, dbo *Option) error {
	db, err := OpenWith(dbo.Host, dbo.Port, dbo.User, dbo.Pass, "")
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbo.DbName))
	return err
}
