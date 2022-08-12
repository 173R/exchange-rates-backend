package repositories

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"gorm.io/gorm"
)

type Currencies struct {
	db *gorm.DB
}

// CreateMany создает множество валют.
func (c *Currencies) CreateMany(currencies []models.Currency) error {
	return c.db.Create(currencies).Error
}

// UpdateById обновляет валюту по её идентификатору.
func (c *Currencies) UpdateById(
	id models.CurrencyId,
	currency *models.Currency,
) error {
	return c.db.Where("id = ?", id).Updates(currency).Error
}

// FindById возвращает валюту по её идентификатору.
func (c *Currencies) FindById(id models.CurrencyId) (*models.Currency, error) {
	var res []models.Currency

	err := c.db.Where("id = ?", id).Limit(1).Find(&res).Error
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], err
}

// FindAll возвращает список всех валют.
func (c *Currencies) FindAll() ([]models.Currency, error) {
	var res []models.Currency

	err := c.db.Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, err
}

// NewCurrencies создает новый экземпляр Currencies.
func NewCurrencies(db *gorm.DB) *Currencies {
	return &Currencies{db: db}
}
