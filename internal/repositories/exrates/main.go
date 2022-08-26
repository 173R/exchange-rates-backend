package exrates

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

// FindLatest возвращает актуальный курс обмена.
func (r *Repository) FindLatest() ([]models.ExchangeRate, error) {
	var res []models.ExchangeRate

	// FIXME: Исправить на gorm.
	if err := r.
		db.
		Raw("select * from exchange_rates t join (select currency_id, MAX(timestamp) as timestamp from exchange_rates group by currency_id) x on x.currency_id = t.currency_id and x.timestamp = t.timestamp").
		Scan(&res).
		Error; err != nil {
		return nil, err
	}
	return res, nil
}

//
//// Create создает новую запись в таблице.
//func (r *Repository) Create(ts time.Time, rates models.RatesJsonb) (*models.ExchangeRate, error) {
//	rec := &models.ExchangeRate{
//		Timestamp: ts,
//		Rates:     rates,
//	}
//
//	err := r.db.Create(rec).Error
//	if err != nil {
//		return nil, err
//	}
//
//	return rec, nil
//}

// NewRepository создает новый экземпляр Repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}
