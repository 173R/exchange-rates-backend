package users

import (
	"context"
	"errors"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/utiltypes"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories/users"
)

type userCurrencies struct {
	rep *users.Repository
}

// FindByUserId возвращает список всех обозреваемых указанным
// пользователем валют.
func (s *userCurrencies) FindByUserId(ctx context.Context, uid models.UserId) (utiltypes.UserObsCurrencies, error) {
	return s.rep.Currencies.FindByUserId(ctx, uid)
}

// Create создает связь между пользователем и валютой.
func (s *userCurrencies) Create(ctx context.Context, uid models.UserId, cid models.CurrencyId) (*models.UserObservedCurrency, error) {
	// Для начала получаем количество отслеживаемых пользователем валют.
	relations, err := s.rep.Currencies.FindByUserId(ctx, uid)
	if err != nil {
		return nil, err
	}

	// Проверяем, отслеживает ли уже пользователь эту валюту.
	for _, r := range relations {
		if r.CurrencyId == cid {
			return r, nil
		}
	}

	// Проверяем, не упёрся ли пользователь в лимит.
	if len(relations) >= 10 {
		return nil, errors.New("observed currencies limit reached")
	}

	return s.rep.Currencies.Create(ctx, uid, cid)
}

// DeleteById удаляет связь пользователя с валютой по идентификатору
// связи.
func (s *userCurrencies) DeleteById(ctx context.Context, id models.UserObservedCurrencyId) (bool, error) {
	return s.rep.Currencies.DeleteById(ctx, id)
}

// DeleteByUserAndCurrencyId удаляет связь пользователя с валютой по его и её
// идентификаторам.
func (s *userCurrencies) DeleteByUserAndCurrencyId(
	ctx context.Context,
	uid models.UserId,
	cid models.CurrencyId,
) (bool, error) {
	return s.rep.Currencies.DeleteByUserAndCurrencyId(ctx, uid, cid)
}
