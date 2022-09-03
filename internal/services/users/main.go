package users

import (
	"context"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories/users"
	"github.com/wolframdeus/exchange-rates-backend/internal/tg"
)

type Service struct {
	rep        *users.Repository
	Currencies *userCurrencies
}

// CreateByTgUid создает стандартного пользователя с указанным Telegram ID.
func (s *Service) CreateByTgUid(ctx context.Context, uid tg.UserId) (*models.User, error) {
	return s.rep.CreateByTgUid(ctx, uid)
}

// FindByTelegramUid возвращает пользователя по его идентификатору Telegram.
func (s *Service) FindByTelegramUid(ctx context.Context, id tg.UserId) (*models.User, error) {
	return s.rep.FindByTelegramUid(ctx, id)
}

// FindById возвращает пользователя по его идентификатору.
func (s *Service) FindById(ctx context.Context, id models.UserId) (*models.User, error) {
	return s.rep.FindById(ctx, id)
}

// SetBaseCurrency устанавливает для пользователя новую базовую валюту.
func (s *Service) SetBaseCurrency(
	ctx context.Context,
	uid models.UserId,
	cid models.CurrencyId,
) (bool, error) {
	return s.rep.SetBase(ctx, uid, cid)
}

// NewService возвращает указатель на новый экземпляр Service.
func NewService(rep *users.Repository) *Service {
	return &Service{
		rep:        rep,
		Currencies: &userCurrencies{rep: rep},
	}
}
