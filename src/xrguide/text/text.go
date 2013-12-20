package text

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
)

var textRefPattern = regexp.MustCompile(`\{(\d+),(\d+)\}`)
var textCommentPattern = regexp.MustCompile(`([^\\]\(.*?[^\\]\))`)

func HasRef(ref string) bool {
	return textRefPattern.MatchString(ref)
}

func ParseTextRef(ref string) (page, text int64, err error) {
	matches := textRefPattern.FindStringSubmatch(ref)
	if len(matches) != 3 {
		err = fmt.Errorf("Invalid text ref format: %s", ref)
		return
	}
	page, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return
	}
	text, err = strconv.ParseInt(matches[2], 10, 64)
	return
}

const selectText = `
SELECT
	text
FROM text_entries
WHERE
	language_id = ?
	AND
	page_id = ?
	AND
	text_id = ?
`

func Get(db *sql.DB, langId, pageId, textId int64) (string, error) {
	var text string
	row := db.QueryRow(selectText, langId, pageId, textId)
	err := row.Scan(&text)
	if err != nil {
		return "", fmt.Errorf("Error getting text: %v", err)
	}
	return text, nil
}

const selectAll = `
SELECT 
	text,
	language_id,
	page_id,
	text_id
FROM text_entries
`

const selectRefTexts = selectAll + `
WHERE
	has_ref = 1
`

const updateText = `
UPDATE text_entries
SET
	text = ?
WHERE
	language_id = ?
	AND
	page_id = ?
	AND
	text_id = ?
`

// Update internal text references
func UpdateReferences(db *sql.DB) error {
	rows, err := db.Query(selectRefTexts)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error retrieving text with refs: %v", err)
	}
	update, err := db.Prepare(updateText)
	if err != nil {
		return fmt.Errorf("Error preparing update stmt: %v", err)
	}
	var text, newText string
	var langId, pageId, textId int64
	var matches []string
	for rows.Next() {
		err = rows.Scan(&text, &langId, &pageId, &textId)
		if err != nil {
			return fmt.Errorf("Error scanning: %v", err)
		}
		err = nil
		newText = textRefPattern.ReplaceAllStringFunc(
			text,
			func(match string) string {
				var repl string
				var pageId, textId int64
				matches = textRefPattern.FindStringSubmatch(match)
				pageId, err = strconv.ParseInt(matches[1], 10, 64)
				if err != nil {
					return ""
				}
				textId, err = strconv.ParseInt(matches[2], 10, 64)
				if err != nil {
					return ""
				}
				repl, err = Get(db, langId, pageId, textId)
				return repl
			},
		)
		if err != nil {
			return fmt.Errorf("Error on replace: %v", err)
		}
		_, err = update.Exec(newText, langId, pageId, textId)
		if err != nil {
			return fmt.Errorf("Error on update: %v", err)
		}
	}
	return nil
}

func UpdateComments(db *sql.DB) error {
	rows, err := db.Query(selectAll)
	if err != nil {
		return fmt.Errorf("Error retrieving text: %v", err)
	}
	//	update, err := db.Prepare(updateText)
	if err != nil {
		return fmt.Errorf("Error preparing update stmt: %v", err)
	}
	var text string
	var langId, pageId, textId int64
	for rows.Next() {
		err = rows.Scan(&text, &langId, &pageId, &textId)
		if err != nil {
			return fmt.Errorf("Error scanning: %v", err)
		}
		if textCommentPattern.MatchString(text) {
			fmt.Println(text)
		}
	}
	return nil
}
