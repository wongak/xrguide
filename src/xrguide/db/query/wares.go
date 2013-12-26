package query

const SelectWaresOverview = `
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
