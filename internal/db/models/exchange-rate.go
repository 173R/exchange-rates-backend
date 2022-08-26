package models

import (
	"time"
)

type ExchangeRate struct {
	// Идентификатор записи.
	Id int64
	// Время обновления данных.
	Timestamp time.Time
	// Идентификатор валюты.
	CurrencyId CurrencyId
	// Абсолютный курс обмена.
	ConvertRate float64
}

// GetExchangeRate возвращает курс обмена валюты base в target.
func (l *ExchangeRate) GetExchangeRate(base CurrencyId, target CurrencyId) (float64, error) {
	return l.GetExchangeRate(base, target)
}
