package users

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"gorm.io/gorm"
)

type Users struct {
	db         *gorm.DB
	Currencies *userCurrencies
}

// FindByTelegramUid возвращает пользователя по его идентификатору Telegram.
func (r *Users) FindByTelegramUid(id int64) (*models.User, error) {
	var res []models.User

	if err := r.db.Where("telegram_uid = ?", id).Limit(1).Find(&res).Error; err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

// SetBaseCurByTgUid обновляет базовую валюту пользователя по его
// идентификатору Telegram.
func (r *Users) SetBaseCurByTgUid(tgUid int64, cid models.CurrencyId) (bool, error) {
	res := r.
		db.
		Model(&models.User{}).
		Where("telegram_uid = ?", tgUid).
		Update("base_currency_id", cid)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

// NewUsers создает новый экземпляр Users.
func NewUsers(db *gorm.DB) *Users {
	return &Users{
		db:         db,
		Currencies: &userCurrencies{db: db},
	}
}
