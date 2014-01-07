package query

const WaresSelectWaresOverview = `
SELECT
	wares.id,
	name_text.text,
	wares.name_raw,
	wares.transport,
	wares.price_average,
	wares.icon
FROM wares
LEFT JOIN text_entries AS name_text ON
	name_text.language_id = ?
	AND
	name_text.page_id = wares.name_page_id
	AND
	name_text.text_id = wares.name_text_id
`

const WaresSelectWare = `
SELECT
	wares.id,
	name_text.text,
	description_text.text,
	wares.name_raw,
	wares.transport,
	wares.specialist,
	wares.size,
	wares.volume,
	wares.price_min,
	wares.price_average,
	wares.price_max,
	wares.container,
	wares.icon
FROM wares
LEFT JOIN text_entries AS name_text ON
	name_text.language_id = ?
	AND
	name_text.page_id = wares.name_page_id
	AND
	name_text.text_id = wares.name_text_id
LEFT JOIN text_entries AS description_text ON
	description_text.language_id = ?
	AND
	description_text.page_id = wares.description_page_id
	AND
	description_text.text_id = wares.description_text_id
WHERE
	wares.id = ?
`

const WaresSelectProductions = `
SELECT
	wares_productions.method,
	wares_productions.time,
	wares_productions.amount,
	name_text.text
FROM
	wares_productions
LEFT JOIN text_entries AS name_text ON
	name_text.language_id = ?
	AND
	name_text.page_id = wares_productions.name_page_id
	AND
	name_text.text_id = wares_productions.name_text_id
WHERE
	wares_productions.ware_id = ?
ORDER BY wares_productions.method, name_text.text
`
