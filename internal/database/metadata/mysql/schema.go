package mysql

var META_TABLE_SCHEMAS = map[string]string{
	"feature": `
		CREATE TABLE feature (
			id				INT 			NOT NULL AUTO_INCREMENT PRIMARY KEY,
			name          	VARCHAR(32) 	NOT	NULL,
			full_name   	VARCHAR(65)		NOT NULL COMMENT "<group_name>:<feature_name>",
			group_id      	INT         	NOT	NULL,
			value_type    	INT         	NOT	NULL COMMENT "data type of feature value",
			description   	VARCHAR(128)	DEFAULT '',
			create_time   	TIMESTAMP    	NOT	NULL DEFAULT CURRENT_TIMESTAMP,
			modify_time   	TIMESTAMP    	NOT	NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE (name)
		);
		`,
	"feature_group": `
		CREATE TABLE feature_group (
			id					INT 			NOT	NULL AUTO_INCREMENT PRIMARY KEY,
			name               	VARCHAR(32) 	NOT	NULL,
			category           	VARCHAR(16) 	NOT	NULL COMMENT "group category: batch, stream",
			entity_id          	INT         	NOT	NULL,
			online_revision_id 	INT         	DEFAULT NULL,
			description        	VARCHAR(64) 	DEFAULT '',
			create_time        	TIMESTAMP   	NOT	NULL DEFAULT CURRENT_TIMESTAMP,
			modify_time			TIMESTAMP   	NOT	NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE (name)
		);
		`,
	"entity": `
		CREATE TABLE entity (
			id				INT 			NOT	NULL AUTO_INCREMENT PRIMARY KEY,
			name        	VARCHAR(32) 	NOT	NULL,
			length      	SMALLINT    	NOT	NULL COMMENT "feature entity value max length",
			description 	VARCHAR(64) 	DEFAULT '',
			create_time 	TIMESTAMP   	NOT	NULL DEFAULT CURRENT_TIMESTAMP,
			modify_time 	TIMESTAMP   	NOT	NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE (name)
		);
		`,
	"feature_group_revision": `
		CREATE TABLE feature_group_revision (
			id				INT 		NOT	NULL AUTO_INCREMENT PRIMARY KEY,
			group_id    	INT         NOT	NULL,
			revision    	BIGINT      NOT	NULL COMMENT "group data point-in-time epoch seconds",
			data_table  	VARCHAR(64) NOT	NULL COMMENT "feature data table name",
			anchored    	BOOLEAN     NOT	NULL,
			description 	VARCHAR(64) DEFAULT '',
			create_time 	TIMESTAMP   NOT	NULL DEFAULT CURRENT_TIMESTAMP,
			modify_time 	TIMESTAMP   NOT	NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE (group_id, revision)
		);
		`,
}

var META_TABLE_FOREIGN_KEYS = map[string]string{
	"feature": `
		ALTER TABLE feature
		ADD FOREIGN KEY (group_id) REFERENCES feature_group(id)
	`,
	"feature_group": `
		ALTER TABLE feature_group
		ADD FOREIGN KEY (entity_id) REFERENCES entity(id),
		ADD FOREIGN KEY (online_revision_id) REFERENCES feature_group_revision(id)
	`,
	"feature_group_revision": `
		ALTER TABLE feature_group_revision
		ADD FOREIGN KEY (group_id) REFERENCES feature_group(id)
	`,
}

var META_VIEW_SCHEMAS = map[string]string{}
