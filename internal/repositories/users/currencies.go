package users

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/utiltypes"
	"gorm.io/gorm"
)

type userCurrencies struct {
	db *gorm.DB
}

// FindByUserId возвращает список всех обозреваемых указанным
// пользователем валют.
func (r *userCurrencies) FindByUserId(uid models.UserId) (utiltypes.UserObsCurrencies, error) {
	var res []*models.UserObservedCurrency

	if err := r.db.Where("user_id = ?", uid).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// Create создает связь между пользователем и валютой.
func (r *userCurrencies) Create(uid models.UserId, cid models.CurrencyId) (*models.UserObservedCurrency, error) {
	value := &models.UserObservedCurrency{
		UserId:     uid,
		CurrencyId: cid,
	}

	err := r.db.Create(value).Error
	if err != nil {
		return nil, err
	}

	return value, nil
}

// DeleteById удаляет связь пользователя с валютой по идентификатору
// связи.
func (r *userCurrencies) DeleteById(id models.UserObservedCurrencyId) (bool, error) {
	res := r.db.Delete(&models.UserObservedCurrency{}, "id = ?", id)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

// DeleteByUserAndCurrencyId удаляет связь пользователя с валютой по его и её
// идентификаторам.
func (r *userCurrencies) DeleteByUserAndCurrencyId(
	uid models.UserId,
	cid models.CurrencyId,
) (bool, error) {
	res := r.
		db.
		Delete(&models.UserObservedCurrency{}, "user_id = ? AND currency_id = ?", uid, cid)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
