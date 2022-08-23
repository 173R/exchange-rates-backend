package jsonb

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/wolframdeus/exchange-rates-backend/internal/language"
)

type TranslationItem struct {
	// Стандартный несклоняемый перевод.
	Default string `json:"default"`
}

type Translation struct {
	Ru *TranslationItem `json:"ru"`
	En *TranslationItem `json:"en"`
}

// Translate возвращает перевод для указанного языка.
func (t *Translation) Translate(lang language.Lang) string {
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

func (t *Translation) Scan(value any) error {
	return ScanTo(value, t)
}

func (t Translation) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// NewTranslation возвращает новый экземпляр Translation.
func NewTranslation(ru string, en string) *Translation {
	return &Translation{
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
