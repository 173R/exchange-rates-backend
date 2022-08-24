package utiltypes

import "github.com/wolframdeus/exchange-rates-backend/internal/db/models"

type UserObsCurrencies []*models.UserObservedCurrency

// CurrencyIds возвращает список идентификаторов всех валют из текущего слайса.
func (s *UserObsCurrencies) CurrencyIds() []models.CurrencyId {
	slice := *s
	res := make([]models.CurrencyId, len(slice))

	for i, item := range slice {
		res[i] = item.CurrencyId
	}
	return res
}
