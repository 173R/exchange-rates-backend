package repositories

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"gorm.io/gorm"
)

type Translations struct {
	db *gorm.DB
}

// CreateMany создает множество переводов.
func (t *Translations) CreateMany(items []models.Translation) error {
	return t.db.Create(items).Error
}

// UpdateById обновляет перевод по его идентификатору.
func (t *Translations) UpdateById(
	id models.TranslationId,
	item *models.Translation,
) error {
	return t.db.Where("id = ?", id).Updates(item).Error
}

// FindAll возвращает список всех переводов.
func (t *Translations) FindAll() ([]models.Translation, error) {
	var res []models.Translation

	if err := t.db.Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// FindById возвращает перевод по его ключу.
func (t *Translations) FindById(id models.TranslationId) (*models.Translation, error) {
	var res []models.Translation

	if err := t.db.Where("id = ?", id).Limit(1).Find(&res).Error; err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

// FindByIds возвращает переводы по их ключам.
func (t *Translations) FindByIds(ids []models.TranslationId) ([]models.Translation, error) {
	var res []models.Translation

	if err := t.db.Where("id IN ?", ids).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// NewTranslations создает ссылку на новый экземпляр Translations.
func NewTranslations(db *gorm.DB) *Translations {
	return &Translations{db: db}
}
