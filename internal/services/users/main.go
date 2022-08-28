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

// SetBaseCurByTgUid обновляет базовую валюту пользователя по его
// идентификатору Telegram.
func (s *Users) SetBaseCurByTgUid(tgUid int64, cid models.CurrencyId) (bool, error) {
	return s.rep.SetBaseCurByTgUid(tgUid, cid)
}

// NewUsers возвращает указатель на новый экземпляр Users.
func NewUsers(rep *users.Users) *Users {
	return &Users{
		rep:        rep,
		Currencies: &userCurrencies{rep: rep},
	}
}
