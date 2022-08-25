package oer

import (
	"time"
)

type latest struct {
	// Дата обновления данных.
	TimestampRaw int64 `mapstructure:"timestamp"`
	// Курс обмена валют.
	Rates map[string]float64 `mapstructure:"rates"`
}

// Timestamp возвращает дату обновления данных в более удобном формате.
func (l *latest) Timestamp() *time.Time {
	t := time.Unix(l.TimestampRaw, 0)
	return &t
}
