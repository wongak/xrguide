package text

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
)

var textRefPattern = regexp.MustCompile(`\{(\d+),(\d+)\}`)

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
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("Error getting text: %v", err)
	}
	return text, nil
}

// Update internal text references
func UpdateReferences(db *sql.DB, langId int64) error {
	return nil
}
