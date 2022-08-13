package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/wolframdeus/exchange-rates-backend/graph/generated"
	"github.com/wolframdeus/exchange-rates-backend/graph/model"
	ctxpkg "github.com/wolframdeus/exchange-rates-backend/internal/context"
)

// Currencies is the resolver for the currencies field.
func (r *queryResolver) Currencies(ctx context.Context) ([]*model.Currency, error) {
	c := ctxpkg.NewGraph(&ctx)

	// Получаем список валют.
	currencies, err := c.Services.Currencies.FindAll()
	if err != nil {
		return nil, err
	}

	fmt.Printf("%+v\n", c.LaunchParams)

	res := make([]*model.Currency, len(currencies))
	for i, c := range currencies {
		res[i] = &model.Currency{
			ID:   string(c.Id),
			Sign: c.Sign,
			// TODO: Добавить перевод.
			Title: c.TitleKey + " (will be translated)",
		}
	}

	return res, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
