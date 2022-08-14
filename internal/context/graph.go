package context

import (
	"context"
	"github.com/wolframdeus/exchange-rates-backend/internal/language"
	"github.com/wolframdeus/exchange-rates-backend/internal/launchparams"
)

type Graph struct {
	// Оригинальный контекст.
	ctx *context.Context
	// TODO: Общие поля с Gin: Services, LaunchParams.
	// Список доступных сервисов.
	Services *Services
	// Список параметров запуска.
	LaunchParams *launchparams.Params
}

// GetLanguage возвращает текущий язык запроса.
func (g *Graph) GetLanguage() language.Lang {
	if g.LaunchParams != nil {
		lang := language.Lang(g.LaunchParams.Language)

		switch lang {
		case language.RU, language.EN:
			return lang
		}
	}
	return language.Default
}

// NewGraph возвращает ссылку на новый экземпляр Graph.
func NewGraph(ctx *context.Context) *Graph {
	c := &Graph{ctx: ctx}
	ctxValue := *ctx

	// Восстанавливаем список сервисов из контекста.
	if srv := getServicesFromContext(ctxValue); srv != nil {
		c.Services = srv
	}

	// Восстанавливаем параметры запуска.
	if params := getLaunchParamsFromContext(ctxValue); params != nil {
		c.LaunchParams = params
	}

	return c
}
