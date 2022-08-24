package users

import (
	"errors"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/utiltypes"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories/users"
)

type userCurrencies struct {
	rep *users.Users
}

// FindByUserId возвращает список всех обозреваемых указанным
// пользователем валют.
func (s *userCurrencies) FindByUserId(uid models.UserId) (utiltypes.UserObsCurrencies, error) {
	return s.rep.Currencies.FindByUserId(uid)
}

// Create создает связь между пользователем и валютой.
func (s *userCurrencies) Create(uid models.UserId, cid models.CurrencyId) (*models.UserObservedCurrency, error) {
	// Для начала получаем количество отслеживаемых пользователем валют.
	relations, err := s.rep.Currencies.FindByUserId(uid)
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

	return s.rep.Currencies.Create(uid, cid)
}

// DeleteById удаляет связь пользователя с валютой по идентификатору
// связи.
func (s *userCurrencies) DeleteById(id models.UserObservedCurrencyId) (bool, error) {
	return s.rep.Currencies.DeleteById(id)
}
