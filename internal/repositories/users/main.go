package users

import (
	"context"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/tg"
	"gorm.io/gorm"
)

type Repository struct {
	db         *gorm.DB
	Currencies *userCurrencies
}

// CreateByTgUid создает стандартного пользователя с указанным Telegram ID.
func (r *Repository) CreateByTgUid(ctx context.Context, uid tg.UserId) (*models.User, error) {
	u := &models.User{
		TelegramUid: uid,
		// TODO: Вынести в константу.
		BaseCurrencyId: "USD",
	}

	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

// FindByTelegramUid возвращает пользователя по его идентификатору Telegram.
func (r *Repository) FindByTelegramUid(ctx context.Context, id tg.UserId) (*models.User, error) {
	var res []models.User

	if err := r.db.WithContext(ctx).Where("telegram_uid = ?", id).Limit(1).Find(&res).Error; err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

// FindById возвращает пользователя по его идентификатору.
func (r *Repository) FindById(ctx context.Context, id models.UserId) (*models.User, error) {
	var res []models.User

	if err := r.db.WithContext(ctx).Where("id = ?", id).Limit(1).Find(&res).Error; err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

// SetBase обновляет базовую валюту пользователя.
func (r *Repository) SetBase(
	ctx context.Context,
	uid models.UserId,
	cid models.CurrencyId,
) (bool, error) {
	res := r.
		db.
		WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", uid).
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
