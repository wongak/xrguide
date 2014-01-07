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
	production_text.text
FROM
	wares_productions
LEFT JOIN text_entries AS production_text ON
	production_text.language_id = ?
	AND
	production_text.page_id = wares_productions.name_page_id
	AND
	production_text.text_id = wares_productions.name_text_id
WHERE
	wares_productions.ware_id = ?
ORDER BY wares_productions.method
`

const WaresSelectProductionWares = `
SELECT
	w.is_primary,
	w.ware,
	name_text.text,
	w.amount
FROM
	wares_productions
INNER JOIN wares_production_wares AS w ON
	w.ware_id = wares_productions.ware_id
	AND
	w.method = wares_productions.method
INNER JOIN wares ON
	w.ware = wares.id
LEFT JOIN text_entries AS name_text ON
	name_text.language_id = ?
	AND
	name_text.page_id = wares.name_page_id
	AND
	name_text.text_id = wares.name_text_id
WHERE
	wares_productions.ware_id = ?
	AND
	wares_productions.method = ?
ORDER BY w.is_primary DESC, name_text.text
`
