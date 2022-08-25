package models

import (
	"fmt"
	"time"
)

type ExchangeRate struct {
	// Идентификатор записи.
	Id int64
	// Время обновления данных.
	Timestamp time.Time
	// Курс обмена валют.
	Rates RatesJsonb
}

// GetExchangeRate возвращает курс обмена валюты base в target.
func (l *ExchangeRate) GetExchangeRate(base CurrencyId, target CurrencyId) (float64, error) {
	return l.GetExchangeRate(base, target)
}

// GetConvertRate возвращает абсолютный курс обмена валюты.
func (l *ExchangeRate) GetConvertRate(cid CurrencyId) (float64, error) {
	rate, ok := l.Rates[cid]
	if ok {
		return rate, nil
	}
	return 0, fmt.Errorf("currency %q not found", cid)
}
