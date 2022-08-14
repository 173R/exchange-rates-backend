package models

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/jsonb"
)

// CurrencyId представляет собой идентификатор валюты.
type CurrencyId string

type Currency struct {
	// Аббревиатура валюты.
	Id CurrencyId
	// Ключ перевода наименования валюты.
	TitleTranslationId TranslationId
	// Информация об изображении.
	Images jsonb.Image
	// Символ валюты.
	Sign string
}
