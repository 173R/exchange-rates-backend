package currencies

import (
	"context"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

// CreateMany создает множество валют.
func (c *Repository) CreateMany(ctx context.Context, currencies []models.Currency) error {
	return c.db.WithContext(ctx).Create(currencies).Error
}

// UpdateById обновляет валюту по её идентификатору.
func (c *Repository) UpdateById(
	ctx context.Context,
	id models.CurrencyId,
	currency *models.Currency,
) error {
	return c.db.WithContext(ctx).Where("id = ?", id).Updates(currency).Error
}

// FindAll возвращает список всех валют.
func (c *Repository) FindAll(ctx context.Context) ([]*models.Currency, error) {
	var res []*models.Currency

	if err := c.db.WithContext(ctx).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// FindById возвращает валюту по её идентификатору.
func (c *Repository) FindById(ctx context.Context, id models.CurrencyId) (*models.Currency, error) {
	var res []models.Currency

	err := c.db.WithContext(ctx).Where("id = ?", id).Limit(1).Find(&res).Error
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], err
}

// FindByIds возвращает валюты по их идентификаторам.
func (c *Repository) FindByIds(ctx context.Context, ids []models.CurrencyId) ([]*models.Currency, error) {
	if len(ids) == 0 {
		return []*models.Currency{}, nil
	}
	var res []*models.Currency

	err := c.db.WithContext(ctx).Where("id IN ?", ids).Limit(len(ids)).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, err
}

// NewRepository создает новый экземпляр Repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}
