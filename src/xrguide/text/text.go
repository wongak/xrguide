package text

import (
	"fmt"
	"regexp"
	"strconv"
)

var textRefPattern = regexp.MustCompile(`^\{(\d+),(\d+)\}$`)

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
