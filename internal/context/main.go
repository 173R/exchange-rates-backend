package context

import (
	"context"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/language"
	"github.com/wolframdeus/exchange-rates-backend/internal/launchparams"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/currencies"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/exrates"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/users"
	"sync"
)

const (
	// Ключ контекста, в котором хранятся сервисы.
	contextKeyServices = "__services"
	// Ключ контекста, в котором хранятся параметры запуска.
	contextKeyLaunchParams = "__launchParams"
	// Ключ контекста, в котором хранится информация о пользователе.
	contextKeyUser = "__user"
	// HeaderLaunchParams - наименование заголовка, в котором хранятся параметры запуска.
	HeaderLaunchParams = "x-launch-params"
)

type Services struct {
	// Сервис для работы с валютами.
	Currencies *currencies.Currencies
	// Сервис для работы с пользователями.
	Users *users.Users
	// Сервис для работы с курсами обменов валют.
	ExchangeRates *exrates.Service
}

// Context представляет собой контекст, который может быть использован как
// в Gin, так и в GraphQL.
type Context struct {
	// Список доступных сервисов.
	Services *Services
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

// GetUser возвращает информацию о текущем пользователе.
func (c *Context) GetUser(ctx *context.Context) (*models.User, error) {
	if c.IsAnonymous() {
		return nil, nil
	}

	// FIXME: Этот код не работает потому что мы работаем с разными экземплярами
	//  Context. Возможно, мьютекс необходимо хранить в самом контексте.
	// Блокируем остальные горутины.
	c.mu.Lock()

	// Не забываем освободить поток.
	defer c.mu.Unlock()

	// Пытаемся получить информацию о пользователе, которую получали раннее.
	cvalue := (*ctx).Value(contextKeyUser)

	// Мы нашли информацию о пользователе.
	if cu, ok := cvalue.(*models.User); ok {
		return cu, nil
	}

	// Мы уже ранее пытались получить эту информацию.
	if _, ok := cvalue.(bool); ok {
		return nil, nil
	}

	// Получаем информацию о пользователе.
	u, err := c.Services.Users.FindByTelegramUid(c.LaunchParams.UserId)
	if err != nil {
		return nil, err
	}

	// Складируем пользователя в текущий контекст либо флаг того, что попытка
	// получения уже проводилась.
	if u == nil {
		*ctx = context.WithValue(*ctx, contextKeyUser, false)
	} else {
		*ctx = context.WithValue(*ctx, contextKeyUser, u)
	}

	return u, nil
}

// Создает новый экземпляр Context.
func newContext(ctx context.Context) *Context {
	c := &Context{}

	// Восстанавливаем список сервисов из контекста.
	if srv := getFromContext[Services](ctx, contextKeyServices); srv != nil {
		c.Services = srv
	}

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
