package services

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories"
)

type Users struct {
	rep *repositories.Users
}

// FindByTelegramUid возвращает пользователя по его идентификатору Telegram.
func (s *Users) FindByTelegramUid(id int64) (*models.User, error) {
	return s.rep.FindByTelegramUid(id)
}

// NewUsers возвращает ссылку на новый экземпляр Users.
func NewUsers(rep *repositories.Users) *Users {
	return &Users{rep: rep}
}
