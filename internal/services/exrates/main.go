package exrates

import (
	"errors"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories/exrates"
	"time"
)

type Service struct {
	rep *exrates.Repository
}

// FindLatest возвращает актуальный курс обмена.
func (s *Service) FindLatest() ([]*models.ExchangeRate, error) {
	return s.rep.FindLatest()
}

// FindLatestByCurrencyId возвращает актуальный курс обмена указанной валюты.
func (s *Service) FindLatestByCurrencyId(cid models.CurrencyId) (*models.ExchangeRate, error) {
	// Получаем список всех свежих курсов обмена.
	latest, err := s.FindLatest()
	if err != nil {
		return nil, err
	}

	// Находим курс обмена по идентификатору валюты.
	for _, r := range latest {
		if r.CurrencyId == cid {
			return r, nil
		}
	}
	return nil, nil
}

// FindPrevDay возвращает последние полученные курсы обмена за предыдущий
// день.
func (s *Service) FindPrevDay() ([]*models.ExchangeRate, error) {
	// Мы приводим текущую дату к началу дня по UTC и используем в дальнейшем
	// запросе.
	ts := time.Now().In(time.UTC).Truncate(24 * time.Hour)

	return s.rep.FindByTimestamp(ts)
}

// FindPrevDayByCurrencyId возвращает курс обмена указанной валюты в
// предыдущий день.
func (s *Service) FindPrevDayByCurrencyId(cid models.CurrencyId) (*models.ExchangeRate, error) {
	// Получаем список всех курсов обмена за предыдущий день.
	rates, err := s.FindPrevDay()
	if err != nil {
		return nil, err
	}

	// Находим курс обмена по идентификатору валюты.
	for _, r := range rates {
		if r.CurrencyId == cid {
			return r, nil
		}
	}
	return nil, nil
}

// FindPrevDayDiff находит абсолютное и процентное отклонение курса валюты
// от предыдущего дня.
func (s *Service) FindPrevDayDiff(cid models.CurrencyId) (float64, float64, error) {
	// Получаем самые свежий курс.
	latest, err := s.FindLatestByCurrencyId(cid)
	if err != nil {
		return 0, 0, err
	}
	if latest == nil {
		return 0, 0, errors.New("latest exchange rate not found")
	}

	// Получаем данные за предыдущий день.
	prevDay, err := s.FindPrevDayByCurrencyId(cid)
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

// NewService возвращает указатель на новый экземпляр Service.
func NewService(rep *exrates.Repository) *Service {
	return &Service{rep: rep}
}
