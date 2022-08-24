package model

type User struct {
	// Список отслеживаемых пользователем валют.
	ObservedCurrencies []*Currency `json:"observedCurrencies"`
	BaseCurrencyId     string      `json:"-"`
	// Базовая валюта пользователя.
	BaseCurrency *Currency `json:"baseCurrency"`
}
