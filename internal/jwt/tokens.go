package jwt

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/language"
)

// UserAccessToken описывает токен, который выдается пользователю для доступа
// к методам API.
type UserAccessToken struct {
	// Идентификатор пользователя.
	Uid models.UserId `mapstructure:"uid" json:"uid"`
	// Язык, используемый пользователем.
	Language language.Lang `mapstructure:"lng" json:"lng"`
}
