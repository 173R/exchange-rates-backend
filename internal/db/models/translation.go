package models

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/wolframdeus/exchange-rates-backend/internal/language"
)

type TranslationItem struct {
	// Стандартный несклоняемый перевод.
	Default string `json:"default"`
}

type TranslationJsonb struct {
	Ru *TranslationItem `json:"ru"`
	En *TranslationItem `json:"en"`
}

// Translate возвращает перевод для указанного языка.
func (t *TranslationJsonb) Translate(lang language.Lang) string {
	switch lang {
	case language.RU:
		return t.Ru.Default
	case language.EN:
		return t.En.Default
	default:
		// FIXME
		return ""
	}
}

func (t *TranslationJsonb) Scan(value any) error {
	return ScanTo(value, t)
}

func (t TranslationJsonb) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// NewTranslationJsonb возвращает новый экземпляр TranslationJsonb.
func NewTranslationJsonb(ru string, en string) *TranslationJsonb {
	return &TranslationJsonb{
		Ru: &TranslationItem{Default: ru},
		En: &TranslationItem{Default: en},
	}
}

/*
Пример корректных данных:
{
	ru: {
		default: "Доллар",
		inclined: ["{count} доллар", "{count} доллара", "{count} долларов"],
	},
	en: { ... }
}
*/
