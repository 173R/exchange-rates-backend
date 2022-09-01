package exrates

import (
	"context"
	"errors"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories/exrates"
)

type Service struct {
	rep   *exrates.Repository
	cache *cache
}

// FindLatestByCurrencyId возвращает актуальный курс обмена указанной валюты.
func (s *Service) FindLatestByCurrencyId(ctx context.Context, cid models.CurrencyId) (*models.ExchangeRate, error) {
	return s.cache.FindLatestById(ctx, cid)
}

// FindPrevDayDiff находит абсолютное и процентное отклонение курса валюты
// от предыдущего дня.
func (s *Service) FindPrevDayDiff(ctx context.Context, cid models.CurrencyId) (float64, float64, error) {
	// Получаем самые свежий курс.
	latest, err := s.FindLatestByCurrencyId(ctx, cid)
	if err != nil {
		return 0, 0, err
	}
	if latest == nil {
		return 0, 0, errors.New("latest exchange rate not found")
	}

	// Получаем данные за предыдущий день.
	prevDay, err := s.findPrevDayByCurrencyId(ctx, cid)
	if err != nil {
		return 0, 0, err
	}
	if prevDay == nil {
		return 0, 0, nil
	}

	// Вычисляем разницу.
	diff := latest.ConvertRate - prevDay.ConvertRate

	return diff, diff / prevDay.ConvertRate * 100, nil
}

// Возвращает курс обмена указанной валюты в предыдущий день.
func (s *Service) findPrevDayByCurrencyId(ctx context.Context, cid models.CurrencyId) (*models.ExchangeRate, error) {
	return s.cache.FindPrevDayById(ctx, cid)
}

// NewService возвращает указатель на новый экземпляр Service.
func NewService(rep *exrates.Repository) *Service {
	return &Service{
		rep:   rep,
		cache: newCache(rep),
	}
}
