package models

type UserId int64
type UserObservedCurrencyId int64

type User struct {
	// Идентификатор пользователя.
	Id UserId
	// Идентификатор пользователя Telegram.
	TelegramUid int64
	// Базовая валюта, выбранная пользователем.
	BaseCurrencyId CurrencyId
}

type UserObservedCurrency struct {
	// Идентификатор связи.
	Id UserObservedCurrencyId
	// Идентификатор пользователя.
	UserId UserId
	// Идентификатор валюты.
	CurrencyId CurrencyId
}
