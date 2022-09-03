package context

import (
	"context"
)

// Context представляет собой контекст, который может быть использован как
// в Gin, так и в GraphQL.
type Context struct {
}

// Извлекает из контекста указанный тип по указанному ключу.
func getFromContext[T interface{}](ctx context.Context, key string) *T {
	if v, ok := ctx.Value(key).(*T); ok {
		return v
	}
	return nil
}
