package currencies

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories/currencies"
)

type Currencies struct {
	cache *cache
}

// FindAll возвращает информацию о всех валютах.
func (c *Currencies) FindAll() ([]*models.Currency, error) {
	return c.cache.FindAll()
}

// FindById возвращает валюту по её идентификатору.
func (c *Currencies) FindById(id models.CurrencyId) (*models.Currency, error) {
	return c.cache.FindById(id)
}

// FindByIds возвращает валюты по их идентификаторам.
func (c *Currencies) FindByIds(ids []models.CurrencyId) ([]*models.Currency, error) {
	return c.cache.FindByIds(ids)
}

// New возвращает ссылку на новый экземпляр Currencies.
func New(rep *currencies.Currencies) *Currencies {
	return &Currencies{cache: newCache(rep)}
}
