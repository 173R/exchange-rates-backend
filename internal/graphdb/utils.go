package graphdb

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/graph/model"
	"github.com/wolframdeus/exchange-rates-backend/internal/language"
)

// CurrencyFromDb создает модель валюты из её модели БД.
func CurrencyFromDb(c *models.Currency, convertRate float64, lang language.Lang) *model.Currency {
	return &model.Currency{
		ID:          string(c.Id),
		Sign:        c.Sign,
		Title:       c.GetTitle(lang),
		ConvertRate: convertRate,
	}
}
