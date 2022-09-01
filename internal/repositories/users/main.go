package users

import (
	"context"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"gorm.io/gorm"
)

type Repository struct {
	db         *gorm.DB
	Currencies *userCurrencies
}

// FindByTelegramUid возвращает пользователя по его идентификатору Telegram.
func (r *Repository) FindByTelegramUid(ctx context.Context, id int64) (*models.User, error) {
	var res []models.User

	if err := r.db.WithContext(ctx).Where("telegram_uid = ?", id).Limit(1).Find(&res).Error; err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

// SetBaseCurByTgUid обновляет базовую валюту пользователя по его
// идентификатору Telegram.
func (r *Repository) SetBaseCurByTgUid(ctx context.Context, tgUid int64, cid models.CurrencyId) (bool, error) {
	res := r.
		db.
		WithContext(ctx).
		Model(&models.User{}).
		Where("telegram_uid = ?", tgUid).
		Update("base_currency_id", cid)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

// NewRepository создает новый экземпляр Repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db:         db,
		Currencies: &userCurrencies{db: db},
	}
}
