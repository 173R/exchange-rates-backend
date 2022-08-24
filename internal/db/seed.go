package db

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/jsonb"
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
	rep := currencies.NewCurrencies(db)

	// Получаем список всех валют.
	items, err := rep.FindAll()
	if err != nil {
		return err
	}

	var toCreate []models.Currency
	var toUpdate []models.Currency
	seed := []models.Currency{
		{
			Id:    "USD",
			Sign:  "$",
			Title: jsonb.NewTranslation("Доллар", "Dollar"),
		},
		{
			Id:    "EUR",
			Sign:  "€",
			Title: jsonb.NewTranslation("Евро", "Euro"),
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
		if err := rep.CreateMany(toCreate); err != nil {
			return err
		}
	}

	// Обновляем существующие валюты.
	if len(toUpdate) > 0 {
		for _, c := range toUpdate {
			if err := rep.UpdateById(c.Id, &c); err != nil {
				return err
			}
		}
	}
	return nil
}
