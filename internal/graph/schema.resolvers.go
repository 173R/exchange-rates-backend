package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	ctxpkg "github.com/wolframdeus/exchange-rates-backend/internal/context"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/graph/generated"
	model2 "github.com/wolframdeus/exchange-rates-backend/internal/graph/model"
	"github.com/wolframdeus/exchange-rates-backend/internal/graphdb"
)

// ConvertRate is the resolver for the convertRate field.
func (r *currencyResolver) ConvertRate(ctx context.Context, obj *model2.Currency) (float64, error) {
	c := ctxpkg.NewGraph(ctx)

	// Получаем курсы обмена валют.
	latest, err := c.Services.ExchangeRates.FindLatest()
	if err != nil {
		return 0, err
	}

	var convertRate *float64
	id := models.CurrencyId(obj.ID)

	for _, rate := range latest {
		if rate.CurrencyId == id {
			convertRate = &rate.ConvertRate
			break
		}
	}

	if convertRate == nil {
		return 0, errors.New("convert rate not found")
	}
	return *convertRate, nil
}

// AddUserObsCurrency is the resolver for the addUserObsCurrency field.
func (r *mutationResolver) AddUserObsCurrency(ctx context.Context, currencyID string) (bool, error) {
	c := ctxpkg.NewGraph(ctx)

	if c.IsAnonymous() {
		return false, errors.New("user not authorized")
	}

	// Получаем информацию о пользователе.
	u, err := c.GetUser(&ctx)
	if err != nil {
		return false, err
	}

	// Создаём связь.
	_, err = c.
		Services.
		Users.
		Currencies.
		Create(u.Id, models.CurrencyId(currencyID))
	if err != nil {
		return false, err
	}

	return true, nil
}

// SetUserBaseCurrency is the resolver for the setUserBaseCurrency field.
func (r *mutationResolver) SetUserBaseCurrency(ctx context.Context, currencyID string) (bool, error) {
	c := ctxpkg.NewGraph(ctx)

	if c.IsAnonymous() {
		return false, errors.New("user not authorized")
	}

	// Обновляем базовую валюту.
	updated, err := c.
		Services.
		Users.
		UpdateBaseCurByTgUid(c.LaunchParams.UserId, models.CurrencyId(currencyID))
	if err != nil {
		return false, err
	}

	return updated, nil
}

// Currencies is the resolver for the currencies field.
func (r *queryResolver) Currencies(ctx context.Context) ([]*model2.Currency, error) {
	c := ctxpkg.NewGraph(ctx)

	// Получаем список валют.
	currencies, err := c.Services.Currencies.FindAll()
	if err != nil {
		return nil, err
	}

	// Получаем язык для того, чтобы перевести наименования валют.
	lang := c.Language()

	res := make([]*model2.Currency, len(currencies))
	for i, cur := range currencies {
		res[i] = graphdb.CurrencyFromDb(cur, 0, lang)
	}

	return res, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context) (*model2.User, error) {
	c := ctxpkg.NewGraph(ctx)

	// Для получения информации о текущем пользователе необходимо быть
	// авторизованным пользователем.
	if c.IsAnonymous() {
		return nil, errors.New("user is not authorized")
	}

	// Находим пользователя по его идентификатору Telegram.
	u, err := c.Services.Users.FindByTelegramUid(c.LaunchParams.UserId)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}

	return &model2.User{BaseCurrencyId: string(u.BaseCurrencyId)}, nil
}

// ObservedCurrencies is the resolver for the observedCurrencies field.
func (r *userResolver) ObservedCurrencies(ctx context.Context, obj *model2.User) ([]*model2.Currency, error) {
	c := ctxpkg.NewGraph(ctx)

	// Получаем информацию о пользователе.
	u, err := c.GetUser(&ctx)
	if err != nil {
		return nil, err
	}

	// Получаем список связей пользователя с валютами.
	relations, err := c.
		Services.
		Users.
		Currencies.
		FindByUserId(u.Id)
	if err != nil {
		return nil, err
	}

	// Находим сами валюты.
	currencies, err := c.Services.Currencies.FindByIds(relations.CurrencyIds())
	if err != nil {
		return nil, err
	}

	// Определяем язык пользователя.
	lang := c.Language()

	// Конвертируем валюты в модели.
	res := make([]*model2.Currency, len(currencies))
	for i, cur := range currencies {
		res[i] = graphdb.CurrencyFromDb(cur, 0, lang)
	}

	return res, nil
}

// BaseCurrency is the resolver for the baseCurrency field.
func (r *userResolver) BaseCurrency(ctx context.Context, obj *model2.User) (*model2.Currency, error) {
	c := ctxpkg.NewGraph(ctx)

	// Находим валюту.
	cur, err := c.Services.Currencies.FindById(models.CurrencyId(obj.BaseCurrencyId))
	if err != nil {
		return nil, err
	}
	if cur == nil {
		return nil, errors.New("base currency not found")
	}

	return graphdb.CurrencyFromDb(cur, 0, c.Language()), nil
}

// Currency returns generated.CurrencyResolver implementation.
func (r *Resolver) Currency() generated.CurrencyResolver { return &currencyResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type currencyResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
