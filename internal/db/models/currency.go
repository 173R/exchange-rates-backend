package models

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/language"
)

// CurrencyId представляет собой идентификатор валюты.
type CurrencyId string

type Currency struct {
	// Аббревиатура валюты.
	Id CurrencyId
	// Переводы наименования валюты.
	Title *TranslationJsonb
	// Информация об изображении.
	Images *ImageJsonb
	// Символ валюты.
	Sign string
}

// GetTitle возвращает заголовок валюты на указанном языке.
func (c *Currency) GetTitle(lang language.Lang) string {
	return c.Title.Translate(lang)
}
