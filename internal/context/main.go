package context

import (
	"context"
	"github.com/wolframdeus/exchange-rates-backend/internal/language"
	"github.com/wolframdeus/exchange-rates-backend/internal/launchparams"
	"github.com/wolframdeus/exchange-rates-backend/internal/services"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/currencies"
)

const (
	// Ключ контекста, в котором хранятся сервисы.
	contextKeyServices = "__services"
	// Ключ контекста, в котором хранятся параметры запуска.
	contextKeyLaunchParams = "__launchParams"
	// HeaderLaunchParams - наименование заголовка, в котором хранятся параметры запуска.
	HeaderLaunchParams = "x-launch-params"
)

type Services struct {
	// Сервис для работы с валютами.
	Currencies *currencies.Currencies
	// Сервис для работы с пользователями.
	Users *services.Users
}

// Context представляет собой контекст, который может быть использован как
// в Gin, так и в GraphQL.
type Context struct {
	// Список доступных сервисов.
	Services *Services
	// Список параметров запуска.
	LaunchParams *launchparams.Params
}

// GetLanguage возвращает текущий язык запроса.
func (c *Context) GetLanguage() language.Lang {
	if c.IsAnonymous() {
		return language.Default
	}
	return c.LaunchParams.Language
}

// IsAnonymous возвращает true в случае, если запрос выполняется анонимно.
func (c *Context) IsAnonymous() bool {
	return c.LaunchParams == nil
}

// NewContext создает новый экземпляр Context.
func NewContext(ctx context.Context) *Context {
	c := &Context{}

	// Восстанавливаем список сервисов из контекста.
	if srv := getServicesFromContext(ctx); srv != nil {
		c.Services = srv
	}

	// Восстанавливаем параметры запуска.
	if params := getLaunchParamsFromContext(ctx); params != nil {
		c.LaunchParams = params
	}

	return c
}

// Извлекает Services из контекста.
func getServicesFromContext(ctx context.Context) *Services {
	return getFromContext[Services](ctx, contextKeyServices)
}

// Извлекает Services из контекста.
func getLaunchParamsFromContext(ctx context.Context) *launchparams.Params {
	return getFromContext[launchparams.Params](ctx, contextKeyLaunchParams)
}

// Извлекает из контекста указанный тип по указанному ключу.
func getFromContext[T interface{}](ctx context.Context, key string) *T {
	if v, ok := ctx.Value(key).(*T); ok {
		return v
	}
	return nil
}
