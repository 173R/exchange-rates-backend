package users

import (
	"context"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/utiltypes"
	"gorm.io/gorm"
)

type userCurrencies struct {
	db *gorm.DB
}

// FindByUserId возвращает список всех обозреваемых указанным
// пользователем валют.
func (r *userCurrencies) FindByUserId(ctx context.Context, uid models.UserId) (utiltypes.UserObsCurrencies, error) {
	var res []*models.UserObservedCurrency

	if err := r.db.WithContext(ctx).Where("user_id = ?", uid).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// Create создает связь между пользователем и валютой.
func (r *userCurrencies) Create(ctx context.Context, uid models.UserId, cid models.CurrencyId) (*models.UserObservedCurrency, error) {
	value := &models.UserObservedCurrency{
		UserId:     uid,
		CurrencyId: cid,
	}

	err := r.db.WithContext(ctx).Create(value).Error
	if err != nil {
		return nil, err
	}

	return value, nil
}

// DeleteById удаляет связь пользователя с валютой по идентификатору
// связи.
func (r *userCurrencies) DeleteById(ctx context.Context, id models.UserObservedCurrencyId) (bool, error) {
	res := r.db.WithContext(ctx).Delete(&models.UserObservedCurrency{}, "id = ?", id)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

// DeleteByUserAndCurrencyId удаляет связь пользователя с валютой по его и её
// идентификаторам.
func (r *userCurrencies) DeleteByUserAndCurrencyId(
	ctx context.Context,
	uid models.UserId,
	cid models.CurrencyId,
) (bool, error) {
	res := r.
		db.
		WithContext(ctx).
		Delete(&models.UserObservedCurrency{}, "user_id = ? AND currency_id = ?", uid, cid)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
