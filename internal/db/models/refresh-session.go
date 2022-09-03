package models

import "time"

type RefreshSessionId int64

type RefreshSession struct {
	// Идентификатор записи.
	Id RefreshSessionId
	// Идентификатор пользователя.
	UserId UserId
	// Токен для обновления токена доступа.
	RefreshToken string
	// Привязанный токен доступа.
	AccessToken string
	// Отпечаток клиента.
	Fingerprint string
	// Дата истечения сессии.
	ExpiresAt time.Time
	// Дата создания сессии.
	CreatedAt time.Time
}
