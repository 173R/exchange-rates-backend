package models

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/jsonb"
	"github.com/wolframdeus/exchange-rates-backend/internal/language"
)

type TranslationId string

type Translation struct {
	// Ключ перевода.
	Id TranslationId
	Ru TranslationJsonb
	En TranslationJsonb
}

// Translate возвращает перевод для указанного языка.
func (t *Translation) Translate(lang language.Lang) string {
	switch lang {
	case language.RU:
		return t.Ru.Default
	case language.EN:
		return t.En.Default
	}
	return string(t.Id)
}

type TranslationJsonb struct {
	// Стандартный несклоняемый перевод.
	Default string `json:"default"`
}

func (v *TranslationJsonb) Scan(value any) error {
	return jsonb.ScanTo(value, v)
}

func (v TranslationJsonb) Value() (driver.Value, error) {
	return json.Marshal(v)
}

/*
Пример корректных данных:
{
	standard: "Доллар",
	inclined: ["{count} доллар", "{count} доллара", "{count} долларов"],
}
*/
