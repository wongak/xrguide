package query

const MacrosSelectStations = `
SELECT
	macros.id,
	name_text.text,
	wares.id,
	ware_name_text.text
FROM macros
INNER JOIN macro_connections ON
	macro_connections.macro_id = macros.id
INNER JOIN macros AS connmacro ON
	connmacro.id = macro_connections.connection_macro_id
INNER JOIN macro_productions ON
	connmacro.id = macro_productions.macro_id
INNER JOIN wares ON
	wares.id = macro_productions.ware_id
LEFT JOIN text_entries AS name_text ON
	name_text.language_id = ?
	AND
	name_text.page_id = macros.name_page_id
	AND
	name_text.text_id = macros.name_text_id
LEFT JOIN text_entries AS ware_name_text ON
	ware_name_text.language_id = ?
	AND
	ware_name_text.page_id = wares.name_page_id
	AND
	ware_name_text.text_id = wares.name_text_id
WHERE
	macros.class = 'station'
GROUP BY macros.id, wares.id
`

const MacrosSelectWareStations = `
SELECT
	macros.id,
	name_text.text
FROM macros
INNER JOIN macro_connections ON
	macro_connections.macro_id = macros.id
INNER JOIN macros AS connmacro ON
	connmacro.id = macro_connections.connection_macro_id
INNER JOIN macro_productions ON
	connmacro.id = macro_productions.macro_id
INNER JOIN wares ON
	wares.id = ?
	AND
	wares.id = macro_productions.ware_id
LEFT JOIN text_entries AS name_text ON
	name_text.language_id = ?
	AND
	name_text.page_id = macros.name_page_id
	AND
	name_text.text_id = macros.name_text_id
WHERE
	macros.class = 'station'
GROUP BY macros.id
`

const MacrosSelectStationProductions = `
SELECT
	macros.id,
	name_text.text,
	wares.id,
	connmacro.id,
	conn_name_text.text,
	macro_connections.mode,
	macro_connections.conn_group,
	macro_connections.sequence,
	macro_connections.stage,
	COUNT(wares.id)
FROM macros
INNER JOIN macro_connections ON
	macro_connections.macro_id = macros.id
INNER JOIN macros AS connmacro ON
	connmacro.id = macro_connections.connection_macro_id
INNER JOIN macro_productions ON
	connmacro.id = macro_productions.macro_id
INNER JOIN wares ON
	wares.id = macro_productions.ware_id
LEFT JOIN text_entries AS name_text ON
	name_text.language_id = ?
	AND
	name_text.page_id = macros.name_page_id
	AND
	name_text.text_id = macros.name_text_id
LEFT JOIN text_entries AS conn_name_text ON
	conn_name_text.language_id = ?
	AND
	conn_name_text.page_id = connmacro.name_page_id
	AND
	conn_name_text.text_id = connmacro.name_text_id
WHERE
	macros.class = 'station'
GROUP BY
	macros.id,
	wares.id,
	macro_connections.sequence,
	macro_connections.stage
`
