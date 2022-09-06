package models

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/language"
	"github.com/wolframdeus/exchange-rates-backend/internal/tg"
)

type UserId int64
type UserObservedCurrencyId int64

type User struct {
	// Идентификатор пользователя.
	Id UserId
	// Идентификатор пользователя Telegram.
	TelegramUid tg.UserId
	// Базовая валюта, выбранная пользователем.
	BaseCurrencyId CurrencyId
	// Предпочитаемый язык пользователя.
	Lang language.Lang
}

type UserObservedCurrency struct {
	// Идентификатор связи.
	Id UserObservedCurrencyId
	// Идентификатор пользователя.
	UserId UserId
	// Идентификатор валюты.
	CurrencyId CurrencyId
}
