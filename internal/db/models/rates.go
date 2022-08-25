package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type RatesJsonb map[CurrencyId]float64

// GetExchangeRate возвращает курс обмена валюты base в target.
func (j *RatesJsonb) GetExchangeRate(base CurrencyId, target CurrencyId) (float64, error) {
	// Получаем курсы обмена для обеих валют.
	baseValue, ok := (*j)[base]
	if !ok {
		return 0, fmt.Errorf("currency %q not found", base)
	}

	targetValue, ok := (*j)[target]
	if !ok {
		return 0, fmt.Errorf("currency %q not found", target)
	}

	return targetValue / baseValue, nil
}

func (j *RatesJsonb) Scan(value any) error {
	return ScanTo(value, j)
}

func (j RatesJsonb) Value() (driver.Value, error) {
	return json.Marshal(j)
}
