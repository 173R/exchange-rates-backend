package utils

import "time"

// TimeISO конвертирует переданное время к строке формата ISO.
func TimeISO(t time.Time) string {
	return t.Format(time.RFC3339)
}
