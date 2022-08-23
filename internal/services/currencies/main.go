package currencies

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories"
)

type Currencies struct {
	cache *cache
}

// FindAll возвращает информацию о всех валютах.
func (c *Currencies) FindAll() ([]models.Currency, error) {
	return c.cache.FindAll()
}

// FindById возвращает валюту по её идентификатору.
func (c *Currencies) FindById(id models.CurrencyId) (*models.Currency, error) {
	return c.cache.FindById(id)
}

// New возвращает ссылку на новый экземпляр Currencies.
func New(rep *repositories.Currencies) *Currencies {
	return &Currencies{cache: newCache(rep)}
}
