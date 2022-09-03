package http

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	ctxpkg "github.com/wolframdeus/exchange-rates-backend/internal/context"
	"github.com/wolframdeus/exchange-rates-backend/internal/graph"
	"github.com/wolframdeus/exchange-rates-backend/internal/graph/generated"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/auth"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/currencies"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/exrates"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/users"
	"strings"
)

// NewGraphQLMiddleware возвращает новый обработчик для GraphQL запросов.
func NewGraphQLMiddleware(
	curSrv *currencies.Service,
	uSrv *users.Service,
	exRatesSrv *exrates.Service,
	authSrv *auth.Service,
) gin.HandlerFunc {
	h := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{
			Resolvers: graph.NewResolver(curSrv, uSrv, exRatesSrv, authSrv),
		}),
	)

	return func(c *gin.Context) {
		// Извлекаем контекст запроса.
		ctx := c.Request.Context()

		// Помещаем в контекст кеш запроса.
		ctx = ctxpkg.ContextWithCache(ctx, uSrv)

		// Помещаем в контекст токен авторизации.
		ctx = context.WithValue(ctx, ctxpkg.KeyAuthToken, deriveAuthToken(c))

		// Переопределяем контекст запроса.
		c.Request = c.Request.WithContext(ctx)

		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Извлекает токен авторизации из запроса.
func deriveAuthToken(c *gin.Context) string {
	h := c.GetHeader("authorization")
	parts := strings.Split(h, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}
