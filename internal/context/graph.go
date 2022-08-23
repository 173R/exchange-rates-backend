package context

import (
	"context"
)

type Graph struct {
	Context
}

// NewGraph возвращает ссылку на новый экземпляр Graph.
func NewGraph(ctx context.Context) *Graph {
	return &Graph{*NewContext(ctx)}
}
