package schema

var MacrosDropMacros = `
DROP TABLE IF EXISTS macros
`

var MacrosCreateMacros = `
CREATE TABLE macros (
	id VARCHAR(128) PRIMARY KEY,
	class VARCHAR(75),
	name_page_id INTEGER(11),
	name_text_id INTEGER(11),
	description_page_id INTEGER(11),
	description_text_id INTEGER(11),
	is_unique TINYINT(1) NOT NULL,
	filename VARCHAR(128)
)
`

var MacrosDropProductions = `
DROP TABLE IF EXISTS macro_productions
`

var MacrosCreateProductions = `
CREATE TABLE macro_productions (
	macro_id VARCHAR(128),
	ware_id VARCHAR(128),
	PRIMARY KEY (macro_id, ware_id)
)
`

var MacrosDropConnections = `
DROP TABLE IF EXISTS macro_connections
`

var MacrosCreateConnections = `
CREATE TABLE macro_connections(
	macro_id VARCHAR(128),
	connection_macro_id VARCHAR(128),
	mode varchar(64),
	conn_group varchar(64),
	sequence char(1),
	stage int
)
`

var MacrosDropIndexes = []string{
	`
DROP INDEX macros_class
	`,
	`
DROP INDEX macros_name_text
	`,
	`
DROP INDEX macros_description_text
	`,
	`
DROP INDEX macro_productions_ware
	`,
	`
DROP INDEX macro_connections_conn
	`,
	`
DROP INDEX macro_connections_build
	`,
}

var MacrosCreateIndexes = []string{
	`
CREATE INDEX macros_class ON macros (class)
	`,
	`
CREATE INDEX macros_name_text ON macros (name_page_id, name_text_id)
	`,
	`
CREATE INDEX macros_description_text ON macros (description_page_id, description_text_id)
	`,
	`
CREATE INDEX macro_productions_ware ON macro_productions (ware_id)
	`,
	`
CREATE INDEX macro_connections_conn ON macro_connections (macro_id, connection_macro_id)
	`,
	`
CREATE INDEX macro_connections_build ON macro_connections (sequence, stage)
	`,
}

var MacrosReset = []*string{
	&MacrosDropMacros,
	&MacrosCreateMacros,
	&MacrosDropProductions,
	&MacrosCreateProductions,
	&MacrosDropConnections,
	&MacrosCreateConnections,
}

const MacrosInsertMacro = `
INSERT INTO macros
(id, class, name_page_id, name_text_id, description_page_id, description_text_id, is_unique, filename)
VALUES
(?, ?, ?, ?, ?, ?, ?, ?)
`

const MacrosInsertProduction = `
INSERT INTO macro_productions 
(macro_id, ware_id)
VALUES
(?, ?)
`

const MacrosInsertConnection = `
INSERT INTO macro_connections
(macro_id, connection_macro_id, mode, conn_group, sequence, stage)
VALUES
(?, ?, ?, ?, ?, ?)
`
