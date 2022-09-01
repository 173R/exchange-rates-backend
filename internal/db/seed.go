package db

import (
	"context"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories/currencies"
	"gorm.io/gorm"
)

// RunSeeds запускает сиды проекта.
func RunSeeds() error {
	db, err := NewByConfig()
	if err != nil {
		return err
	}

	// Выполняем все сиды.
	if err := seedCurrencies(db); err != nil {
		return err
	}
	return nil
}

// Запускает сиды, связанные с валютами.
func seedCurrencies(db *gorm.DB) error {
	rep := currencies.NewRepository(db)
	ctx := context.TODO()

	// Получаем список всех валют.
	items, err := rep.FindAll(ctx)
	if err != nil {
		return err
	}

	var toCreate []models.Currency
	var toUpdate []models.Currency
	seed := []models.Currency{
		{
			Id:    "USD",
			Sign:  "$",
			Title: models.NewTranslationJsonb("Доллар", "Dollar"),
		},
		{
			Id:    "EUR",
			Sign:  "€",
			Title: models.NewTranslationJsonb("Евро", "Euro"),
		},
		{
			Id:    "RUB",
			Sign:  "₽",
			Title: models.NewTranslationJsonb("Рубль", "Ruble"),
			Images: &models.ImageJsonb{
				Set: []models.ImageJsonbSetItem{
					{
						Width:  100,
						Height: 200,
						Url:    "https://privet-zriteli.ru/image.png",
						Scale:  1,
					},
				},
			},
		},
		{
			Id:    "GBP",
			Sign:  "£",
			Title: models.NewTranslationJsonb("Фунт Стерлингов", "Pound Sterling"),
		},
	}

	for _, sItem := range seed {
		found := false

		for _, item := range items {
			if sItem.Id == item.Id {
				toUpdate = append(toUpdate, sItem)
				found = true
				break
			}
		}

		if !found {
			toCreate = append(toCreate, sItem)
		}
	}

	// Создаем отсутствующие валюты.
	if len(toCreate) > 0 {
		if err := rep.CreateMany(ctx, toCreate); err != nil {
			return err
		}
	}

	// Обновляем существующие валюты.
	if len(toUpdate) > 0 {
		for _, c := range toUpdate {
			if err := rep.UpdateById(ctx, c.Id, &c); err != nil {
				return err
			}
		}
	}
	return nil
}
