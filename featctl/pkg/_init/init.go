package _init

import (
	"context"
	"fmt"
	"log"

	"github.com/onestore-ai/onestore/pkg/database"
)

type Option struct {
	DBOption database.Option
}

func createDatabase(ctx context.Context, dbo *database.Option) error {
	db, err := database.OpenWith(dbo.Host, dbo.Port, dbo.User, dbo.Pass, "")
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbo.DbName))
	return err
}

func createFeatureRevisionTable(ctx context.Context, db *database.DB) error {
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
	_, err := db.ExecContext(ctx, schema)
	return err
}

func createFeatureConfigTable(ctx context.Context, db *database.DB) error {
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
	_, err := db.ExecContext(ctx, schema)
	return err
}

func Init(ctx context.Context, option *Option) {
	if err := createDatabase(ctx, &option.DBOption); err != nil {
		log.Fatalf("failed creating database: %v", err)
	}

	db, err := database.Open(&option.DBOption)
	if err != nil {
		log.Fatalf("failed connecting feature store: %v", err)
	}
	defer db.Close()

	if err := createFeatureConfigTable(ctx, db); err != nil {
		log.Fatalf("failed initializing feature store: %v", err)
	}

	if err := createFeatureRevisionTable(ctx, db); err != nil {
		log.Fatalf("failed initializing feature store: %v", err)
	}

	log.Println("succeeded.")
}
