package schema

var TextDropLanguages = `
DROP TABLE IF EXISTS languages;
`
var TextCreateLanguages = `
CREATE TABLE languages (
	id INTEGER PRIMARY KEY,
	name VARCHAR(75) UNIQUE
)
`

var TextDropEntries = `
DROP TABLE IF EXISTS text_entries;
`

var TextCreateEntries = `
CREATE TABLE text_entries (
	language_id INTEGER,
	page_id INTEGER,
	text_id INTEGER,
	text TEXT,
	PRIMARY KEY (language_id, page_id, text_id),
	FOREIGN KEY (language_id) REFERENCES languages(id) ON DELETE RESTRICT ON UPDATE CASCADE
)
`

var TextDefaultLanguages = `
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
`

var TextReset = []*string{
	&TextDropLanguages,
	&TextCreateLanguages,
	&TextDropEntries,
	&TextCreateEntries,
	&TextDefaultLanguages,
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
