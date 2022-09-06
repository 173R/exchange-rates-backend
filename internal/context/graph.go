package context

import (
	"context"
	"errors"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/jwt"
	"github.com/wolframdeus/exchange-rates-backend/internal/language"
)

const (
	KeyAuthToken = "__auth-token"
)

type Graph struct {
	cache *Cache
	token *jwt.UserAccessToken
}

// Language возвращает _известный_ язык запроса.
func (w *Graph) Language() language.Lang {
	if w.token == nil || !w.token.Language.Known() {
		return language.Default
	}
	return w.token.Language
}

// User возвращает информацию о пользователе, совершающий запрос.
func (w *Graph) User(ctx context.Context) (*models.User, error) {
	if w.token == nil {
		return nil, errors.New("user is not authorized")
	}
	return w.cache.GetUser(ctx, w.token.Uid)
}

// UserId возвращает идентификатор пользователя.
func (w *Graph) UserId() (models.UserId, error) {
	if w.token == nil {
		return 0, errors.New("user is not authorized")
	}
	return w.token.Uid, nil
}

// NewGraph возвращает указатель на новый экземпляр Graph.
func NewGraph(ctx context.Context) (*Graph, error) {
	// Извлекаем кеш запроса.
	cache, err := CacheFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Парсим токен доступа пользователя.
	var ut *jwt.UserAccessToken
	token, ok := ctx.Value(KeyAuthToken).(string)

	if ok && len(token) > 0 {
		if utoken, err := jwt.DecodeUserAccessToken(token); err == nil {
			ut = utoken
		}
	}

	return &Graph{cache: cache, token: ut}, nil
}

// MustNewGraph является опасной функцией, которая выбросит панику в
// случае, если при вызове NewGraph произошла ошибка.
func MustNewGraph(ctx context.Context) *Graph {
	w, err := NewGraph(ctx)
	if err != nil {
		panic(err)
	}
	return w
}