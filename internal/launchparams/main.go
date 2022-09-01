package launchparams

import (
	"encoding/json"
	"errors"
	"github.com/wolframdeus/exchange-rates-backend/internal/language"
)

type Params struct {
	// Локаль пользователя в сыром виде.
	LanguageRaw string
	// Локально пользователя для использования в проекте.
	Language language.Lang
	// Идентификатор пользователя.
	UserId int64
}

type paramsJson struct {
	Language string `json:"language"`
	UserId   *int64 `json:"user_id"`
}

func (p *Params) UnmarshalJSON(data []byte) error {
	var jp paramsJson
	if err := json.Unmarshal(data, &jp); err != nil {
		return err
	}

	// Идентификатор пользователя является обязательным.
	if jp.UserId == nil {
		return errors.New("user_id is missing")
	}

	// Определяем локаль пользователя.
	langRaw := jp.Language
	lang := language.Lang(langRaw)

	switch lang {
	case language.RU, language.EN:
		break
	default:
		lang = language.Default
	}

	*p = Params{
		UserId:      *jp.UserId,
		LanguageRaw: langRaw,
		Language:    lang,
	}
	return nil
}

// Derive извлекает параметры запуска из какой-либо строки.
func Derive(value string) (*Params, error) {
	var j Params

	if err := json.Unmarshal([]byte(value), &j); err != nil {
		return nil, err
	}
	return &j, nil
}
