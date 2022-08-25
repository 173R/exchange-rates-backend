package exrates

import (
	"errors"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	db *gorm.DB
}

// FindLatest возвращает актуальный курс обмена.
func (r *Repository) FindLatest() (*models.ExchangeRate, error) {
	var res []models.ExchangeRate

	if err := r.db.Order("timestamp DESC").Limit(1).Find(&res).Error; err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("no exchange rates records found")
	}
	return &res[0], nil
}

// Create создает новую запись в таблице.
func (r *Repository) Create(ts time.Time, rates models.RatesJsonb) (*models.ExchangeRate, error) {
	rec := &models.ExchangeRate{
		Timestamp: ts,
		Rates:     rates,
	}

	err := r.db.Create(rec).Error
	if err != nil {
		return nil, err
	}

	return rec, nil
}

// NewRepository создает новый экземпляр Repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}
