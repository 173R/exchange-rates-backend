package http

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"github.com/wolframdeus/exchange-rates-backend/graph"
	"github.com/wolframdeus/exchange-rates-backend/graph/generated"
)

// NewGraphQLMiddleware возвращает новый обработчик для GraphQL запросов.
func NewGraphQLMiddleware() gin.HandlerFunc {
	h := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}),
	)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
