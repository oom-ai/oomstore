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
	"feature_entity": `
		CREATE TABLE feature_entity (
			id			SERIAL2 PRIMARY KEY,
			name		VARCHAR(32) NOT NULL UNIQUE,
			length		SMALLINT    NOT NULL,

			description	VARCHAR(64) DEFAULT '',
			create_time	TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modify_time	TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		COMMENT ON COLUMN feature_entity.length IS 'feature entity value max length';
		`,
	"feature_group": `
		CREATE TABLE feature_group (
			id					SERIAL2 PRIMARY KEY,
			name				VARCHAR(32) NOT NULL UNIQUE,
			entity_id			INT2 NOT NULL,
			online_revision_id	INT	DEFAULT NULL,
			category			VARCHAR(16) NOT NULL,

			description			VARCHAR(64) DEFAULT '',
			create_time			TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modify_time			TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		COMMENT ON COLUMN feature_group.online_revision_id IS 'group online point-in-time epoch seconds';
		COMMENT ON COLUMN feature_group.category   IS 'group category: batch, stream ...';
		`,
	"feature": `
		CREATE TABLE feature (
			id				SERIAL2 PRIMARY KEY,
			name			VARCHAR(32) NOT NULL UNIQUE,
			group_id		INT2 NOT NULL,
			db_value_type	VARCHAR(32) NOT NULL,
			value_type		VARCHAR(16) NOT NULL,

			description		VARCHAR(128) DEFAULT '',
			create_time		TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modify_time		TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		COMMENT ON COLUMN feature.value_type    IS 'data type of feature value';
		COMMENT ON COLUMN feature.db_value_type IS 'database data type of feature value';
		`,
	"feature_group_revision": `
		CREATE TABLE feature_group_revision (
			id			SERIAL PRIMARY KEY,
			group_id	INT2 NOT NULL,
			group_name  VARCHAR(32) NOT NULL,
			revision	BIGINT      NOT NULL,
			data_table	VARCHAR(64) NOT NULL,

			description	VARCHAR(64) DEFAULT '',
			create_time	TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modify_time	TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

			UNIQUE(group_id, revision)
		);
		COMMENT ON COLUMN feature_group_revision.revision   IS 'group data point-in-time epoch seconds';
		COMMENT ON COLUMN feature_group_revision.data_table IS 'feature data table name';
		`,
}

var META_VIEW_SCHEMAS = map[string]string{
	"rich_feature": `
        	CREATE VIEW rich_feature AS
			SELECT
				tmp3.*,
				feature_entity.name AS entity_name
			FROM
			(SELECT
				tmp2.*,
				fgr2.revision AS online_revision
			FROM
			(SELECT
				tmp.*,
				fgr.revision AS offline_revision,
				fgr.data_table AS offline_data_table
			FROM
			(SELECT
				f.*,
				fg.entity_id, fg.category, fg.online_revision_id, fg.name AS group_name
			FROM feature AS f
			LEFT JOIN feature_group AS fg
			ON f.group_id = fg.id) AS tmp
			LEFT JOIN feature_group_revision AS fgr
			ON
				tmp.group_id = fgr.group_id AND
				fgr.revision = (
					SELECT MAX(revision)
					FROM feature_group_revision
					WHERE feature_group_revision.group_id = tmp.group_id
    			)) AS tmp2
			LEFT JOIN feature_group_revision AS fgr2
			ON tmp2.online_revision_id = fgr2.id) AS tmp3
			LEFT JOIN feature_entity
			ON tmp3.entity_id = feature_entity.id;

	`,
	"rich_feature_group": `
        	CREATE VIEW rich_feature_group AS
			SELECT
				tmp.*,
				feature_entity.name AS entity_name
			FROM
			(SELECT
				fg_tmp.*,
				fgr2.revision AS offline_revision,
				fgr2.data_table AS offline_data_table
			FROM
			(SELECT
				fg.id, fg.name, fg.entity_id, fg.category, fg.online_revision_id, fg.description, fg.create_time, fg.modify_time,
				fgr.revision AS online_revision
			FROM feature_group AS fg
			LEFT JOIN feature_group_revision AS fgr
			ON fg.online_revision_id = fgr.id) AS fg_tmp
			LEFT JOIN feature_group_revision AS fgr2
			ON
				fg_tmp.id = fgr2.group_id AND
				fgr2.revision = (
					SELECT MAX(revision)
					FROM feature_group_revision
					WHERE feature_group_revision.group_id = fg_tmp.id
				)) AS tmp
			LEFT JOIN feature_entity
			ON tmp.entity_id = feature_entity.id;
	`,
}
