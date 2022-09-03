package context

import (
	"context"
	"errors"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/users"
	"sync"
)

type cacheKey[T interface{}] struct {
	// Последнее полученное значение.
	Value T
	// Факт попытки получения этого значения.
	Tried bool
	// Мьютекс, который контролирует доступ к значению.
	mu sync.Mutex
}

type Cache struct {
	userSrv *users.Service
	// Последний полученный пользователь.
	user cacheKey[*models.User]
}

// GetUser возвращает информацию о пользователе по его идентификатору.
func (c *Cache) GetUser(ctx context.Context, uid models.UserId) (*models.User, error) {
	// Блокируем и разблокируем мьютекс связанный с пользователем.
	c.user.mu.Lock()
	defer c.user.mu.Unlock()

	// Возвращаем кешированный результат.
	if c.user.Tried {
		return c.user.Value, nil
	}

	// Получаем информацию о пользователе и кешируем её.
	u, err := c.userSrv.FindById(ctx, uid)
	if err != nil {
		return nil, err
	}

	c.user.Tried = true
	c.user.Value = u

	return u, nil
}

// ContextWithCache возвращает новый контекст с помещенным в него кешем GraphQL.
func ContextWithCache(ctx context.Context, uSrv *users.Service) context.Context {
	return context.WithValue(ctx, "__cache", &Cache{userSrv: uSrv})
}

// CacheFromContext извлекает из контекста кеш GraphQL.
func CacheFromContext(ctx context.Context) (*Cache, error) {
	v, ok := ctx.Value("__cache").(*Cache)
	if !ok {
		return nil, errors.New("graphql cache not found or has incorrect value")
	}
	return v, nil
}
