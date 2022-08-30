package model

type User struct {
	// Идентификатор пользователя.
	Id string `json:"id"`
	// Сырой идентификатор пользователя.
	IdRaw int64 `json:"-"`
	// Список отслеживаемых пользователем валют.
	ObservedCurrencies []*Currency `json:"observedCurrencies"`
	// Идентификатор базовой валюты пользователя.
	BaseCurrencyId string `json:"baseCurrencyId"`
	// Базовая валюта пользователя.
	BaseCurrency *Currency `json:"baseCurrency"`
}
