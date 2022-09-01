package context

import (
	"context"
	"github.com/wolframdeus/exchange-rates-backend/internal/language"
	"github.com/wolframdeus/exchange-rates-backend/internal/launchparams"
	"sync"
)

const (
	// Ключ контекста, в котором хранятся параметры запуска.
	contextKeyLaunchParams = "__launchParams"
	// HeaderLaunchParams - наименование заголовка, в котором хранятся параметры запуска.
	HeaderLaunchParams = "x-launch-params"
)

// Context представляет собой контекст, который может быть использован как
// в Gin, так и в GraphQL.
type Context struct {
	// Список параметров запуска.
	LaunchParams *launchparams.Params
	mu           sync.Mutex
}

// Language возвращает текущий язык запроса.
func (c *Context) Language() language.Lang {
	if c.IsAnonymous() {
		return language.Default
	}
	return c.LaunchParams.Language
}

// IsAnonymous возвращает true в случае, если запрос выполняется анонимно.
func (c *Context) IsAnonymous() bool {
	return c.LaunchParams == nil
}

// Создает новый экземпляр Context.
func newContext(ctx context.Context) *Context {
	c := &Context{}

	// Восстанавливаем параметры запуска.
	if params := getFromContext[launchparams.Params](ctx, contextKeyLaunchParams); params != nil {
		c.LaunchParams = params
	}

	return c
}

// Извлекает из контекста указанный тип по указанному ключу.
func getFromContext[T interface{}](ctx context.Context, key string) *T {
	if v, ok := ctx.Value(key).(*T); ok {
		return v
	}
	return nil
}
