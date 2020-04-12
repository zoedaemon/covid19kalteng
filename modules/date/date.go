package date

import (
	"errors"
	"time"
)

const LAYOUT_DATE_SIMPLE = "2006-01-02"

//ParseSimple date YY-MM-DD
func ParseSimple(date string) (time.Time, error) {

	if len(date) > 0 {
		//parse with default layout format
		t, err := time.Parse(LAYOUT_DATE_SIMPLE, date)
		if err != nil {
			return time.Time{}, err
		}
		return t, nil
	}

	return time.Time{}, errors.New("nil")
}
