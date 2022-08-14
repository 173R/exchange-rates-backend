package services

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories"
)

type Translations struct {
	rep   *repositories.Translations
	cache *translationsCache
}

// FindByIds возвращает переводы по их идентификаторам.
func (t *Translations) FindByIds(ids []models.TranslationId) ([]models.Translation, error) {
	return t.cache.FindByIds(ids)
}

// NewTranslations создает ссылку на новый экземпляр Translations.
func NewTranslations(rep *repositories.Translations) *Translations {
	return &Translations{
		rep:   rep,
		cache: &translationsCache{rep: rep},
	}
}
