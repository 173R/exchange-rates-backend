package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/wolframdeus/exchange-rates-backend/configs"
	"time"
)

type Client struct {
	*redis.Client
}

// InvalidateUserAccessToken инвалидирует токен пользователя.
func (c *Client) InvalidateUserAccessToken(
	ctx context.Context,
	token string,
	ttl time.Duration,
) error {
	return c.Set(ctx, getUATKey(token), true, ttl).Err()
}

// InvalidateUserAccessTokens инвалидирует множество токенов пользователей.
// В качестве значения передается карта, в которой ключ является токеном,
// а значение - срок его действия.
func (c *Client) InvalidateUserAccessTokens(
	ctx context.Context,
	tokens map[string]time.Duration,
) error {
	if len(tokens) == 0 {
		return nil
	}

	_, err := c.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for t, ttl := range tokens {
			pipe.Set(ctx, getUATKey(t), true, ttl)
		}
		return nil
	})
	return err
}

// IsUserAccessTokenValid возвращает true в случае, если переданный токен
// пользователя валиден.
func (c *Client) IsUserAccessTokenValid(
	ctx context.Context,
	token string,
) (bool, error) {
	_, err := c.Get(ctx, getUATKey(token)).Result()
	if err != nil {
		if err == redis.Nil {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

// New создает новый экземпляр клиента Redis.
func New(host string, port uint, pass string) *Client {
	return &Client{
		redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", host, port),
			Password: pass,
			DB:       0,
		}),
	}
}

// NewWithConfig создает новый экземпляр клиента Redis из конфигурации
// проекта.
func NewWithConfig() *Client {
	return New(configs.Redis.Host, configs.Redis.Port, configs.Redis.Pass)
}

// Возвращает наименование ключа Redis соответствующее невалидному токену
// доступа пользователя.
func getUATKey(token string) string {
	return fmt.Sprintf("inv-uat-%s", token)
}
