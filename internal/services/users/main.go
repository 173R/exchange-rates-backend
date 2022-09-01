package users

import (
	"context"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories/users"
)

type Service struct {
	rep        *users.Repository
	Currencies *userCurrencies
}

// FindByTelegramUid возвращает пользователя по его идентификатору Telegram.
func (s *Service) FindByTelegramUid(ctx context.Context, id int64) (*models.User, error) {
	return s.rep.FindByTelegramUid(ctx, id)
}

// SetBaseCurByTgUid обновляет базовую валюту пользователя по его
// идентификатору Telegram.
func (s *Service) SetBaseCurByTgUid(ctx context.Context, tgUid int64, cid models.CurrencyId) (bool, error) {
	return s.rep.SetBaseCurByTgUid(ctx, tgUid, cid)
}

// NewService возвращает указатель на новый экземпляр Service.
func NewService(rep *users.Repository) *Service {
	return &Service{
		rep:        rep,
		Currencies: &userCurrencies{rep: rep},
	}
}
