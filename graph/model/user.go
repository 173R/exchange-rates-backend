package model

type User struct {
	// Список отслеживаемых пользователем валют.
	ObservedCurrencies []*Currency `json:"observedCurrencies"`
	// Идентификатор базовой валюты пользователя.
	BaseCurrencyId string `json:"baseCurrencyId"`
	// Базовая валюта пользователя.
	BaseCurrency *Currency `json:"baseCurrency"`
}
