package mysql

var META_TABLE_SCHEMAS = map[string]string{
	"feature": `
		CREATE TABLE feature (
			id				INT 			NOT NULL AUTO_INCREMENT PRIMARY KEY,
			name          	VARCHAR(32) 	NOT	NULL,
			group_id      	INT         	NOT	NULL,
			value_type    	INT         	NOT	NULL COMMENT "data type of feature value",
			description   	VARCHAR(128)	DEFAULT '',
			create_time   	TIMESTAMP    	NOT	NULL DEFAULT CURRENT_TIMESTAMP,
			modify_time   	TIMESTAMP    	NOT	NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE (group_id, name)
		);
		`,
	"feature_group": `
		CREATE TABLE feature_group (
			id					INT 			NOT	NULL AUTO_INCREMENT PRIMARY KEY,
			name               	VARCHAR(32) 	NOT	NULL,
			category           	VARCHAR(16) 	NOT	NULL COMMENT "group category: batch, stream",
			entity_id          	INT         	NOT	NULL,
			snapshot_interval  	INT		        DEFAULT 0,
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
			snapshot_table  VARCHAR(64) NOT	NULL COMMENT "batch & streaming feature snapshot table name",
			cdc_table    	VARCHAR(64) NOT	NULL DEFAULT '' COMMENT "streaming feature cdc table name",
			anchored    	BOOLEAN     NOT	NULL,
			description 	VARCHAR(64) DEFAULT '',
			create_time 	TIMESTAMP   NOT	NULL DEFAULT CURRENT_TIMESTAMP,
			modify_time 	TIMESTAMP   NOT	NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE (group_id, revision)
		);
		`,
}

var META_TABLE_FOREIGN_KEYS = []string{
	"ALTER TABLE feature ADD FOREIGN KEY (group_id) REFERENCES feature_group(id)",
	// On TiDB, if not explicitly given a name, the following 2 foreign key constraints will share the same name
	// So we given them different names to avoid 'ERROR 1826: Duplicate foreign key constraint name'
	"ALTER TABLE feature_group ADD CONSTRAINT FK_feature_group_entity_id FOREIGN KEY (entity_id) REFERENCES entity(id)",
	"ALTER TABLE feature_group ADD CONSTRAINT FK_feature_group_online_revision_id FOREIGN KEY (online_revision_id) REFERENCES feature_group_revision(id)",
	"ALTER TABLE feature_group_revision ADD FOREIGN KEY (group_id) REFERENCES feature_group(id)",
}

var META_VIEW_SCHEMAS = map[string]string{}
