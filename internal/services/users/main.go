package users

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories/users"
)

type Users struct {
	rep        *users.Users
	Currencies *userCurrencies
}

// FindByTelegramUid возвращает пользователя по его идентификатору Telegram.
func (s *Users) FindByTelegramUid(id int64) (*models.User, error) {
	return s.rep.FindByTelegramUid(id)
}

// UpdateBaseCurByTgUid обновляет базовую валюту пользователя по его
// идентификатору Telegram.
func (s *Users) UpdateBaseCurByTgUid(tgUid int64, cid models.CurrencyId) (bool, error) {
	return s.rep.UpdateBaseCurByTgUid(tgUid, cid)
}

// NewUsers возвращает указатель на новый экземпляр Users.
func NewUsers(rep *users.Users) *Users {
	return &Users{
		rep:        rep,
		Currencies: &userCurrencies{rep: rep},
	}
}
