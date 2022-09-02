package jwt

import "github.com/wolframdeus/exchange-rates-backend/internal/db/models"

// UserAccessToken описывает токен, который выдается пользователю для доступа
// к методам API.
type UserAccessToken struct {
	// Идентификатор пользователя.
	Uid models.UserId `mapstructure:"uid" json:"uid"`
}
