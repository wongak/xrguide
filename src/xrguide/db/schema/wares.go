package schema

var WaresDropWares = `
DROP TABLE IF EXISTS wares
`

var WaresCreateWares = `
CREATE TABLE wares (
	id VARCHAR(128) PRIMARY KEY,
	name_page_id INTEGER(11),
	name_text_id INTEGER(11),
	description_page_id INTEGER(11),
	description_text_id INTEGER(11),
	name_raw VARCHAR(255),
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

var WaresDropProductionWares = `
DROP TABLE IF EXISTS wares_production_wares
`

var WaresCreateProductionWares = `
CREATE TABLE wares_production_wares (
	ware_id VARCHAR(128) NOT NULL,
	method VARCHAR(64) NOT NULL,
	is_primary TINYINT(1),
	ware VARCHAR(128) NOT NULL,
	amount INT,
	FOREIGN KEY (ware_id, method) REFERENCES wares_productions(ware_id, method) ON DELETE RESTRICT ON UPDATE CASCADE,
	FOREIGN KEY (ware) REFERENCES wares(id) ON DELETE RESTRICT ON UPDATE CASCADE
)
`

var WaresDropProductionEffects = `
DROP TABLE IF EXISTS wares_production_effects
`

var WaresCreateProductionEffects = `
CREATE TABLE wares_production_effects (
	ware_id VARCHAR(128) NOT NULL,
	method VARCHAR(64) NOT NULL,
	type VARCHAR(128),
	product FLOAT,
	FOREIGN KEY (ware_id, method) REFERENCES wares_productions(ware_id, method) ON DELETE RESTRICT ON UPDATE CASCADE
)
`

var WaresDropIndexes = []string{
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

var WaresCreateIndexes = []string{
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

var WaresReset = []*string{
	&WaresDropWares,
	&WaresCreateWares,
	&WaresDropProductions,
	&WaresCreateProductions,
	&WaresDropProductionWares,
	&WaresCreateProductionWares,
}

var WareInsertWare = `
INSERT INTO wares
(
id, 
name_page_id, 
name_text_id,
description_page_id,
description_text_id,
name_raw,
transport,
specialist,
size,
volume ,
price_min ,
price_average ,
price_max 
)
VALUES
(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`
