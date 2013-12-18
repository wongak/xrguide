package schema

var TextReset = []string{
	`
DROP TABLE IF EXISTS languages;
		`,
	`
CREATE TABLE languages (
	id INTEGER PRIMARY KEY ASC,
	name TEXT UNIQUE
)
		`,
	`
DROP TABLE IF EXISTS text_entries;
		`,
	`
CREATE TABLE text_entries (
	language_id INTEGER,
	page_id INTEGER,
	text_id INTEGER,
	text TEXT,
	PRIMARY KEY (language_id, page_id, text_id ASC),
	FOREIGN KEY (language_id) REFERENCES languages(id) ON DELETE RESTRICT ON UPDATE CASCADE
)
		`,
	`
INSERT INTO languages
(id, name)
VALUES
(7, 'Russian'),
(33, 'French'),
(34, 'Spanish'),
(39, 'Italian'),
(44, 'English'),
(49, 'German'),
(86, 'Chinese (traditional)'),
(88, 'Chinese (simplified)')
		`,
}

var TextDeletePage string = `
DELETE FROM text_entries
WHERE
language_id = ?
AND
page_id = ?
`
var TextInsert string = `
INSERT INTO text_entries
(language_id, page_id, text_id, text)
VALUES
(?, ?, ?, ?)
`
