package sqlite

var TRIGGER_TEMPLATE = `
CREATE TRIGGER {{TABLE_NAME}}_update_modify_time
AFTER UPDATE ON {{TABLE_NAME}}
 BEGIN
  update {{TABLE_NAME}} SET modify_time = datetime('now') WHERE id = NEW.id;
 END;
`

var META_TABLE_SCHEMAS = map[string]string{
	"feature": `
		CREATE TABLE feature (
			id				INTEGER 		NOT NULL PRIMARY KEY AUTOINCREMENT,
			name          	VARCHAR(32) 	NOT	NULL,
			full_name   	VARCHAR(65)		NOT NULL,
			group_id      	INT         	NOT	NULL,
			value_type    	INT          	NOT	NULL,
			description   	VARCHAR(128)	DEFAULT '',
			create_time   	TIMESTAMP    	NOT	NULL DEFAULT CURRENT_TIMESTAMP,
			modify_time   	TIMESTAMP    	NOT	NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (name),
			FOREIGN KEY (group_id) REFERENCES feature_group(id)
		);
		`,
	"feature_group": `
		CREATE TABLE feature_group (
			id					INTEGER 		NOT	NULL PRIMARY KEY AUTOINCREMENT,
			name               	VARCHAR(32) 	NOT	NULL,
			category           	VARCHAR(16) 	NOT	NULL,
			entity_id          	INT         	NOT	NULL,
			online_revision_id 	INT         	DEFAULT NULL,
			description        	VARCHAR(64) 	DEFAULT '',
			create_time        	TIMESTAMP   	NOT	NULL DEFAULT CURRENT_TIMESTAMP,
			modify_time			TIMESTAMP   	NOT	NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (name),
			FOREIGN KEY (entity_id) REFERENCES entity(id),
			FOREIGN KEY (online_revision_id) REFERENCES feature_group_revision(id)
		);
		`,
	"entity": `
		CREATE TABLE entity (
			id				INTEGER 		NOT	NULL PRIMARY KEY AUTOINCREMENT,
			name        	VARCHAR(32) 	NOT	NULL,
			length      	SMALLINT    	NOT	NULL,
			description 	VARCHAR(64) 	DEFAULT '',
			create_time 	TIMESTAMP   	NOT	NULL DEFAULT CURRENT_TIMESTAMP,
			modify_time 	TIMESTAMP   	NOT	NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (name)
		);
		`,
	"feature_group_revision": `
		CREATE TABLE feature_group_revision (
			id				INTEGER 	NOT	NULL PRIMARY KEY AUTOINCREMENT,
			group_id    	INT         NOT	NULL,
			revision    	BIGINT      NOT	NULL,
			data_table  	VARCHAR(64) NOT	NULL,
			anchored    	BOOLEAN     NOT	NULL,
			description 	VARCHAR(64) DEFAULT '',
			create_time 	TIMESTAMP   NOT	NULL DEFAULT CURRENT_TIMESTAMP,
			modify_time 	TIMESTAMP   NOT	NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (group_id, revision),
			FOREIGN KEY (group_id) REFERENCES feature_group(id)
		);
		`,
}

var META_VIEW_SCHEMAS = map[string]string{}
