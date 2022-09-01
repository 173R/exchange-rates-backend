package currencies

import (
	"context"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories/currencies"
)

type Service struct {
	cache *cache
}

// FindAll возвращает информацию о всех валютах.
func (c *Service) FindAll(ctx context.Context) ([]*models.Currency, error) {
	return c.cache.FindAll(ctx)
}

// FindById возвращает валюту по её идентификатору.
func (c *Service) FindById(ctx context.Context, id models.CurrencyId) (*models.Currency, error) {
	return c.cache.FindById(ctx, id)
}

// FindByIds возвращает валюты по их идентификаторам.
func (c *Service) FindByIds(ctx context.Context, ids []models.CurrencyId) ([]*models.Currency, error) {
	return c.cache.FindByIds(ctx, ids)
}

// NewService возвращает указатель на новый экземпляр Service.
func NewService(rep *currencies.Repository) *Service {
	return &Service{cache: newCache(rep)}
}
