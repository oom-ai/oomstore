package database

var META_TABLE_SCHEMAS = map[string]string{
	"feature": `
		CREATE TABLE feature (
			name        VARCHAR(32) NOT NULL COMMENT 'feature name',
			group_name  VARCHAR(32) NOT NULL COMMENT 'feature group name',
			value_type  VARCHAR(16) NOT NULL COMMENT 'sql data type of feature value',

			description VARCHAR(128) DEFAULT NULL COMMENT 'feature description',
			create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
			modify_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'modify time',
			PRIMARY KEY pk(name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`,
	"feature_group": `
		CREATE TABLE feature_group (
			name        VARCHAR(32) NOT     NULL COMMENT 'group name',
			entity_name VARCHAR(32) NOT     NULL COMMENT 'group entity field name',
			revision    BIGINT      DEFAULT NULL COMMENT 'group online point-in-time epoch seconds',
			category    VARCHAR(16) NOT     NULL COMMENT 'group category: batch / stream',
			data_table  VARCHAR(64) DEFAULT NULL COMMENT 'feature data table name',

			description VARCHAR(64) DEFAULT NULL COMMENT 'group description',
			create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
			modify_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'modify time',
			PRIMARY KEY pk(name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`,
	"feature_entity": `
		CREATE TABLE feature_entity (
			name    VARCHAR(32) NOT NULL COMMENT 'feature entity name',
			length	SMALLINT    NOT NULL COMMENT 'feature entity value max length',

			description VARCHAR(64) DEFAULT NULL COMMENT 'feature entity description',
			create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
			modify_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'modify time',
			PRIMARY KEY pk(name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`,
	"feature_group_revision": `
		CREATE TABLE feature_group_revision (
			group_name  VARCHAR(32) NOT NULL COMMENT 'group name',
			revision    BIGINT      NOT NULL COMMENT 'group data point-in-time epoch seconds',
			data_table  VARCHAR(64) NOT NULL COMMENT 'feature data table name',

			description VARCHAR(64) DEFAULT NULL COMMENT 'group description',
			create_time TIMESTAMP   NOT     NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
			modify_time TIMESTAMP   NOT     NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'modify time',
			PRIMARY KEY pk(group_name, revision)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`,
}

var META_VIEW_SCHEMAS = map[string]string{
	"rich_feature": `
        	CREATE VIEW rich_feature AS
			SELECT
				f.name, f.group_name, f.value_type, f.description, f.create_time, f.modify_time,
				fg.entity_name, fg.category, fg.revision, fg.data_table
			FROM feature AS f
			LEFT JOIN feature_group AS fg
			ON f.group_name = fg.name;
	`,
}
