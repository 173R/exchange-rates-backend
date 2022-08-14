package context

import (
	"context"
	"github.com/wolframdeus/exchange-rates-backend/internal/launchparams"
	"github.com/wolframdeus/exchange-rates-backend/internal/services"
	services2 "github.com/wolframdeus/exchange-rates-backend/internal/services/translations"
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
	Currencies *services.Currencies
	// Сервис для работы с переводами.
	Translations *services2.Translations
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
