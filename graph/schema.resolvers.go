package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"github.com/wolframdeus/exchange-rates-backend/graph/generated"
	"github.com/wolframdeus/exchange-rates-backend/graph/model"
	ctxpkg "github.com/wolframdeus/exchange-rates-backend/internal/context"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
)

// Currencies is the resolver for the currencies field.
func (r *queryResolver) Currencies(ctx context.Context) ([]*model.Currency, error) {
	c := ctxpkg.NewGraph(ctx)

	// Получаем список валют.
	currencies, err := c.Services.Currencies.FindAll()
	if err != nil {
		return nil, err
	}

	// Получаем язык для того, чтобы перевести наименования валют.
	lang := c.GetLanguage()

	res := make([]*model.Currency, len(currencies))
	for i, cur := range currencies {
		res[i] = &model.Currency{
			ID:    string(cur.Id),
			Sign:  cur.Sign,
			Title: cur.GetTitle(lang),
		}
	}

	return res, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	c := ctxpkg.NewGraph(ctx)

	// Для получения информации о текущем пользователе необходимо быть
	// авторизованным пользователем.
	if c.IsAnonymous() {
		return nil, errors.New("user is not authorized")
	}

	u, err := c.Services.Users.FindByTelegramUid(c.LaunchParams.UserId)
	if err != nil {
		return nil, err
	}

	return &model.User{
		BaseCurrencyId: string(u.BaseCurrencyId),
		//ObservedCurrencies: []*model.Currency{},
	}, nil
}

// ObservedCurrencies is the resolver for the observedCurrencies field.
func (r *userResolver) ObservedCurrencies(ctx context.Context, obj *model.User) ([]*model.Currency, error) {
	// FIXME
	return []*model.Currency{}, nil
}

// BaseCurrency is the resolver for the baseCurrency field.
func (r *userResolver) BaseCurrency(ctx context.Context, obj *model.User) (*model.Currency, error) {
	c := ctxpkg.NewGraph(ctx)

	// Находим валюту.
	cur, err := c.Services.Currencies.FindById(models.CurrencyId(obj.BaseCurrencyId))
	if err != nil {
		return nil, err
	}
	if cur == nil {
		return nil, errors.New("base currency not found")
	}

	return &model.Currency{
		ID:    string(cur.Id),
		Title: cur.GetTitle(c.GetLanguage()),
		Sign:  cur.Sign,
	}, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
