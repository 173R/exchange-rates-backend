package exrates

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories/exrates"
)

type Service struct {
	rep *exrates.Repository
}

// FindLatest возвращает актуальный курс обмена.
func (s *Service) FindLatest() (*models.ExchangeRate, error) {
	return s.rep.FindLatest()
}

// NewService возвращает указатель на новый экземпляр Service.
func NewService(rep *exrates.Repository) *Service {
	return &Service{rep: rep}
}
