package models

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/jsonb"
)

type Currency struct {
	// Аббревиатура валюты.
	Id string
	// Ключ перевода наименования валюты.
	TitleKey string
	// Информация об изображении.
	Images jsonb.Image
	// Символ валюты.
	Sign string
}
