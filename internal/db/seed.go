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
	if err := seedCurrencies(db); err != nil {
		return err
	}
	return nil
}

// Запускает сиды, связанные с валютами.
func seedCurrencies(db *gorm.DB) error {
	rep := repositories.NewCurrencies(db)

	// Получаем список всех валют.
	currencies, err := rep.FindAll()
	if err != nil {
		return err
	}

	var toCreate []models.Currency
	var toUpdate []models.Currency
	seed := []models.Currency{
		{
			Id:       "USD",
			Sign:     "$",
			TitleKey: "usd_title",
		},
		{
			Id:       "EUR",
			Sign:     "€",
			TitleKey: "euro_title",
		},
	}

	for _, sc := range seed {
		found := false

		for _, c := range currencies {
			if sc.Id == c.Id {
				toUpdate = append(toUpdate, sc)
				found = true
				break
			}
		}

		if !found {
			toCreate = append(toCreate, sc)
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
