package graphdb

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/graph/model"
	"github.com/wolframdeus/exchange-rates-backend/internal/language"
)

// CurrencyFromDb создает модель валюты из её модели БД.
func CurrencyFromDb(c *models.Currency, lang language.Lang) *model.Currency {
	return &model.Currency{
		ID:     string(c.Id),
		Sign:   c.Sign,
		Title:  c.GetTitle(lang),
		Images: ImagesFromDb(c.Images),
	}
}

// ImagesFromDb создает слайс graphql-моделей Image из их модели БД.
func ImagesFromDb(img *models.ImageJsonb) []*model.Image {
	if img == nil {
		return []*model.Image{}
	}
	res := make([]*model.Image, len(img.Set))

	for i, v := range img.Set {
		res[i] = &model.Image{
			Width:  int(v.Width),
			Height: int(v.Height),
			URL:    v.Url,
			Scale:  int(v.Scale),
		}
	}
	return res
}
