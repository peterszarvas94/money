package utils

import (
	"time"
)

func ConvertToTime(date string) (time.Time, error) {
	layout := time.RFC3339Nano
	t, err := time.Parse(layout, date)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
