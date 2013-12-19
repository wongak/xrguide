package schema

var WaresDropWares = `
DROP TABLE IF EXISTS wares
`

var WaresCreateWares = `
CREATE TABLE wares (
	id VARCHAR(128) PRIMARY KEY,
	name_page_id INTEGER(11) NOT NULL,
	name_text_id INTEGER(11) NOT NULL,
	description_page_id INTEGER(11) NOT NULL,
	description_text_id INTEGER(11) NOT NULL,
	transport VARCHAR(64),
	specialist VARCHAR(128),
	size VARCHAR(64),
	volume INTEGER,
	price_min INTEGER,
	price_average INTEGER,
	price_max INTEGER,
	FOREIGN KEY (name_page_id, name_text_id) REFERENCES text_entries(page_id, text_id) ON DELETE RESTRICT ON UPDATE CASCADE,
	FOREIGN KEY (description_page_id, description_text_id) REFERENCES text_entries(page_id, text_id) ON DELETE RESTRICT ON UPDATE CASCADE
)
`

var WaresDropWaresIndexes = []string{
	`
DROP INDEX wares_transport
	`,
	`
DROP INDEX wares_specialist
	`,
	`
DROP INDEX wares_size
	`,
}

var WaresCreateWaresIndexes = []string{
	`
CREATE INDEX wares_transport ON wares (transport)
	`,
	`
CREATE INDEX wares_specialist ON wares (specialist)
	`,
	`
CREATE INDEX wares_size ON wares (size)
	`,
}

var WaresDropProductions = `
DROP TABLE IF EXISTS wares_productions
`

var WaresCreateProductions = `
CREATE TABLE wares_productions (
	ware_id VARCHAR(128),
	method VARCHAR(64),
	time INTEGER,
	amount INTEGER,
	name_page_id INTEGER NOT NULL,
	name_text_id INTEGER NOT NULL,
	PRIMARY KEY (ware_id, method),
	FOREIGN KEY (ware_id) REFERENCES wares(id) ON DELETE RESTRICT ON UPDATE CASCADE,
	FOREIGN KEY (name_page_id, name_text_id) REFERENCES text_entries(page_id, text_id) ON DELETE RESTRICT ON UPDATE CASCADE
)
`

var WaresReset = []*string{
	&WaresDropWares,
	&WaresCreateWares,
}
