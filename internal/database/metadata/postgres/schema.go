package postgres

var DB_FUNCTIONS = []string{
	`
	CREATE OR REPLACE FUNCTION update_modify_time() RETURNS trigger AS $$
		BEGIN
			NEW.modify_time = NOW();
			RETURN NEW;
		END;
	$$ LANGUAGE plpgsql;
	`,
}

var TRIGGER_TEMPLATE = `
	CREATE TRIGGER {{TABLE_NAME}}_update_modify_time
	BEFORE UPDATE ON {{TABLE_NAME}}
	FOR EACH ROW
	EXECUTE PROCEDURE update_modify_time();
`

var META_TABLE_SCHEMAS = map[string]string{
	"feature": `
		CREATE TABLE feature (
			id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			name        	VARCHAR(32)  	NOT NULL,
			group_id    	INT          	NOT NULL,
			value_type  	INT  		   	NOT NULL,
			description 	VARCHAR(128) 	DEFAULT '',
			create_time   	TIMESTAMP    	NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modify_time   	TIMESTAMP    	NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (group_id, name)
		);
		COMMENT ON COLUMN feature.value_type IS 'data type of feature value';
		`,
	"feature_group": `
		CREATE TABLE feature_group (
			id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			name               VARCHAR(32) NOT     NULL,
			category           VARCHAR(16) NOT     NULL,
			entity_id          INT         NOT     NULL,
			snapshot_interval  INT         DEFAULT 0,
			online_revision_id INT         DEFAULT NULL,
			description        VARCHAR(64) DEFAULT '',
			create_time        TIMESTAMP   NOT     NULL DEFAULT CURRENT_TIMESTAMP,
			modify_time        TIMESTAMP   NOT     NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (name)
		);
		COMMENT ON COLUMN feature_group.category IS 'group category: batch, stream ...';
		`,
	"entity": `
		CREATE TABLE entity (
			id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			name        VARCHAR(32) NOT     NULL,
			description VARCHAR(64) DEFAULT '',
			create_time TIMESTAMP   NOT     NULL DEFAULT CURRENT_TIMESTAMP,
			modify_time TIMESTAMP   NOT     NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (name)
		);
		`,
	"feature_group_revision": `
		CREATE TABLE feature_group_revision (
			id    	    	INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			group_id    	INT         NOT NULL,
			revision    	BIGINT      NOT NULL,
			snapshot_table  VARCHAR(64) NOT	NULL,
			cdc_table    	VARCHAR(64) NOT	NULL DEFAULT '',
			anchored    	BOOLEAN     NOT NULL,
			description 	VARCHAR(64) DEFAULT '',
			create_time 	TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modify_time 	TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (group_id, revision)
		);
		COMMENT ON COLUMN feature_group_revision.revision   IS 'group data point-in-time epoch seconds';
		COMMENT ON COLUMN feature_group_revision.snapshot_table IS 'batch & streaming feature snapshot table name';
		COMMENT ON COLUMN feature_group_revision.cdc_table IS 'streaming feature cdc table name';
		`,
}

var META_TABLE_FOREIGN_KEYS = map[string]string{
	"feature": `
		ALTER TABLE feature
		ADD CONSTRAINT fk_group
			FOREIGN KEY(group_id)
			REFERENCES feature_group(id)
	`,
	"feature_group": `
		ALTER TABLE feature_group
		ADD CONSTRAINT fk_entity
				FOREIGN KEY(entity_id)
				REFERENCES entity(id),
		ADD CONSTRAINT fk_online_revision
				FOREIGN KEY(online_revision_id)
				REFERENCES feature_group_revision(id)
	`,
	"feature_group_revision": `
		ALTER TABLE feature_group_revision
		ADD CONSTRAINT fk_group
			FOREIGN KEY(group_id)
			REFERENCES feature_group(id)
	`,
}

var META_VIEW_SCHEMAS = map[string]string{}
