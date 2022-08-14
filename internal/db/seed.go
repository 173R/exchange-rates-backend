package db

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories"
	"gorm.io/gorm"
)

// RunSeeds запускает сиды проекта.
func RunSeeds() error {
	db, err := NewByConfig()
	if err != nil {
		return err
	}

	// Выполняем все сиды.
	if err := seedTranslations(db); err != nil {
		return err
	}
	if err := seedCurrencies(db); err != nil {
		return err
	}
	return nil
}

// Запускает сиды, связанные с валютами.
func seedCurrencies(db *gorm.DB) error {
	rep := repositories.NewCurrencies(db)

	// Получаем список всех валют.
	items, err := rep.FindAll()
	if err != nil {
		return err
	}

	var toCreate []models.Currency
	var toUpdate []models.Currency
	seed := []models.Currency{
		{
			Id:                 "USD",
			Sign:               "$",
			TitleTranslationId: "usd_title",
		},
		{
			Id:                 "EUR",
			Sign:               "€",
			TitleTranslationId: "euro_title",
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

// Запускает сиды, связанные с переводами.
func seedTranslations(db *gorm.DB) error {
	rep := repositories.NewTranslations(db)

	// Получаем список всех валют.
	items, err := rep.FindAll()
	if err != nil {
		return err
	}

	var toCreate []models.Translation
	var toUpdate []models.Translation
	seed := []models.Translation{
		{
			Id: "usd_title",
			Ru: models.TranslationJsonb{Default: "доллар"},
			En: models.TranslationJsonb{Default: "dollar"},
		},
		{
			Id: "euro_title",
			Ru: models.TranslationJsonb{Default: "евро"},
			En: models.TranslationJsonb{Default: "euro"},
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

	// Создаем отсутствующие переводы.
	if len(toCreate) > 0 {
		if err := rep.CreateMany(toCreate); err != nil {
			return err
		}
	}

	// Обновляем существующие переводы.
	if len(toUpdate) > 0 {
		for _, c := range toUpdate {
			if err := rep.UpdateById(c.Id, &c); err != nil {
				return err
			}
		}
	}
	return nil
}
