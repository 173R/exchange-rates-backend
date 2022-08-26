package context

import (
	"context"
)

type Graph struct {
	Context
}

// NewGraph возвращает указатель на новый экземпляр Graph.
func NewGraph(ctx context.Context) *Graph {
	return &Graph{*newContext(ctx)}
}
