package repositories

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"gorm.io/gorm"
)

type Users struct {
	db *gorm.DB
}

// FindByTelegramUid возвращает пользователя по его идентификатору Telegram.
func (r *Users) FindByTelegramUid(id int64) (*models.User, error) {
	var res []models.User

	err := r.db.Where("telegram_uid = ?", id).Limit(1).Find(&res).Error
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], err
}

// NewUsers создает новый экземпляр Users.
func NewUsers(db *gorm.DB) *Users {
	return &Users{db: db}
}
