package jwt

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/language"
	"time"
)

const (
	tokenTypeUserAccessToken  tokenType = "user_access_token"
	tokenTypeUserRefreshToken tokenType = "user_refresh_token"
)

// Описывает список типов токена, которые может генерировать пакет.
type tokenType string

// UserAccessToken описывает токен, который выдается пользователю для доступа
// к методам API.
type UserAccessToken struct {
	// Идентификатор пользователя.
	Uid models.UserId `mapstructure:"uid" json:"uid"`
	// Язык, используемый пользователем.
	Language language.Lang `mapstructure:"lng" json:"lng"`
	// Unix-время выдачи токена.
	IssuedAtRaw int64 `mapstructure:"iat" json:"iat"`
	// Unix-время истечения токена.
	ExpiresAtRaw int64 `mapstructure:"exp" json:"exp"`
}

// IssuedAt возвращает дату выдачи токена.
func (t *UserAccessToken) IssuedAt() time.Time {
	return time.Unix(t.IssuedAtRaw, 0)
}

// ExpiresAt возвращает дату окончания действия токена.
func (t *UserAccessToken) ExpiresAt() time.Time {
	return time.Unix(t.ExpiresAtRaw, 0)
}

// ExpiresIn возвращает оставшуюся длительность действия токена.
func (t *UserAccessToken) ExpiresIn() time.Duration {
	diff := t.ExpiresAt().Sub(time.Now())
	if diff < 0 {
		return 0
	}
	return diff
}
