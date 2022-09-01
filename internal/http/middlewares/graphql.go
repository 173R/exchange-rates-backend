package http

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"github.com/wolframdeus/exchange-rates-backend/internal/graph"
	"github.com/wolframdeus/exchange-rates-backend/internal/graph/generated"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/currencies"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/exrates"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/users"
)

// NewGraphQLMiddleware возвращает новый обработчик для GraphQL запросов.
func NewGraphQLMiddleware(
	curSrv *currencies.Service,
	uSrv *users.Service,
	exRatesSrv *exrates.Service,
) gin.HandlerFunc {
	h := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{
			Resolvers: graph.NewResolver(curSrv, uSrv, exRatesSrv),
		}),
	)

	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(graph.ContextWithCache(c.Request.Context(), uSrv))
		h.ServeHTTP(c.Writer, c.Request)
	}
}
