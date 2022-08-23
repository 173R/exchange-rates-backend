package model

type User struct {
	BaseCurrencyId string `json:"-"`
	// Список отслеживаемых пользователем валют.
	ObservedCurrencies []*Currency `json:"observedCurrencies"`
	// Базовая валюта пользователя.
	BaseCurrency *Currency `json:"baseCurrency"`
}
