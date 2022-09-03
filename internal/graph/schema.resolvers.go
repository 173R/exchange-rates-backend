package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"strconv"

	ctxpkg "github.com/wolframdeus/exchange-rates-backend/internal/context"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/graph/generated"
	"github.com/wolframdeus/exchange-rates-backend/internal/graph/model"
	"github.com/wolframdeus/exchange-rates-backend/internal/graphdb"
	"github.com/wolframdeus/exchange-rates-backend/internal/utils"
)

// ConvertRate is the resolver for the convertRate field.
func (r *currencyResolver) ConvertRate(ctx context.Context, obj *model.Currency) (*model.CurrencyConvertRate, error) {
	cid := models.CurrencyId(obj.ID)

	// Получаем курс обмена этой валюты.
	rate, err := r.Services.ExchangeRates.FindLatestByCurrencyId(ctx, cid)
	if err != nil {
		return nil, err
	}
	if rate == nil {
		return nil, nil
	}

	return &model.CurrencyConvertRate{
		Rate:      rate.ConvertRate,
		UpdatedAt: utils.TimeISO(rate.Timestamp),
	}, nil
}

// Diff is the resolver for the diff field.
func (r *currencyResolver) Diff(ctx context.Context, obj *model.Currency) (*model.CurrencyDiff, error) {
	// Находим изменение валюты.
	adiff, pdiff, err := r.
		Services.
		ExchangeRates.
		FindPrevDayDiff(ctx, models.CurrencyId(obj.ID))
	if err != nil {
		return nil, err
	}

	return &model.CurrencyDiff{
		Absolute: adiff,
		Percents: pdiff,
	}, nil
}

// AuthenticateTg is the resolver for the authenticateTg field.
func (r *mutationResolver) AuthenticateTg(ctx context.Context, initData string, fp string) (*model.AuthResult, error) {
	// Аутентифицируем пользователя.
	res, err := r.Services.Auth.AuthenticateTg(ctx, initData, fp)
	if err != nil {
		return nil, err
	}

	return graphdb.AuthResultFromResult(res), nil
}

// AddUserObsCurrency is the resolver for the addUserObsCurrency field.
func (r *mutationResolver) AddUserObsCurrency(ctx context.Context, currencyID string) (bool, error) {
	c, err := ctxpkg.NewGraph(ctx)
	if err != nil {
		return false, err
	}

	// Получаем ID пользователя.
	uid, err := c.UserId()
	if err != nil {
		return false, err
	}

	// Создаём связь.
	_, err = r.
		Services.
		Users.
		Currencies.
		Create(ctx, uid, models.CurrencyId(currencyID))
	if err != nil {
		return false, err
	}

	return true, nil
}

// RemoveUserObsCurrency is the resolver for the removeUserObsCurrency field.
func (r *mutationResolver) RemoveUserObsCurrency(ctx context.Context, currencyID string) (bool, error) {
	c, err := ctxpkg.NewGraph(ctx)
	if err != nil {
		return false, err
	}

	// Получаем идентификатор пользователя.
	uid, err := c.UserId()
	if err != nil {
		return false, err
	}

	// Удаляем связь пользователя с валютой.
	return r.Services.Users.Currencies.DeleteByUserAndCurrencyId(ctx, uid, models.CurrencyId(currencyID))
}

// RefreshSession is the resolver for the refreshSession field.
func (r *mutationResolver) RefreshSession(ctx context.Context, refreshToken string, fp string) (*model.AuthResult, error) {
	// Обновляем сессию.
	res, err := r.Services.Auth.RefreshSession(ctx, refreshToken, fp)
	if err != nil {
		return nil, err
	}

	return graphdb.AuthResultFromResult(res), nil
}

// SetUserBaseCurrency is the resolver for the setUserBaseCurrency field.
func (r *mutationResolver) SetUserBaseCurrency(ctx context.Context, currencyID string) (bool, error) {
	c, err := ctxpkg.NewGraph(ctx)
	if err != nil {
		return false, err
	}

	// Получаем идентификатор пользователя.
	uid, err := c.UserId()
	if err != nil {
		return false, err
	}

	// Обновляем базовую валюту.
	updated, err := r.
		Services.
		Users.
		SetBaseCurrency(ctx, uid, models.CurrencyId(currencyID))
	if err != nil {
		return false, err
	}

	return updated, nil
}

// Currencies is the resolver for the currencies field.
func (r *queryResolver) Currencies(ctx context.Context) ([]*model.Currency, error) {
	c, err := ctxpkg.NewGraph(ctx)
	if err != nil {
		return nil, err
	}

	// Получаем список валют.
	currencies, err := r.Services.Currencies.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Получаем язык для того, чтобы перевести наименования валют.
	lang := c.Language()

	res := make([]*model.Currency, len(currencies))
	for i, cur := range currencies {
		res[i] = graphdb.CurrencyFromDb(cur, lang)
	}

	return res, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	c, err := ctxpkg.NewGraph(ctx)
	if err != nil {
		return nil, err
	}

	// Находим пользователя по его идентификатору Telegram.
	u, err := c.User(ctx)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}

	return &model.User{
		Id:             strconv.FormatUint(uint64(u.Id), 10),
		IdRaw:          int64(u.Id),
		BaseCurrencyId: string(u.BaseCurrencyId),
	}, nil
}

// ObservedCurrencies is the resolver for the observedCurrencies field.
func (r *userResolver) ObservedCurrencies(ctx context.Context, obj *model.User) ([]*model.Currency, error) {
	c, err := ctxpkg.NewGraph(ctx)
	if err != nil {
		return nil, err
	}

	// Получаем список связей пользователя с валютами.
	relations, err := r.
		Services.
		Users.
		Currencies.
		FindByUserId(ctx, models.UserId(obj.IdRaw))
	if err != nil {
		return nil, err
	}

	// Находим сами валюты.
	currencies, err := r.Services.Currencies.FindByIds(ctx, relations.CurrencyIds())
	if err != nil {
		return nil, err
	}

	// Определяем язык пользователя.
	lang := c.Language()

	// Конвертируем валюты в модели.
	res := make([]*model.Currency, len(currencies))
	for i, cur := range currencies {
		res[i] = graphdb.CurrencyFromDb(cur, lang)
	}

	return res, nil
}

// BaseCurrency is the resolver for the baseCurrency field.
func (r *userResolver) BaseCurrency(ctx context.Context, obj *model.User) (*model.Currency, error) {
	c, err := ctxpkg.NewGraph(ctx)
	if err != nil {
		return nil, err
	}

	// Находим валюту.
	cur, err := r.Services.Currencies.FindById(ctx, models.CurrencyId(obj.BaseCurrencyId))
	if err != nil {
		return nil, err
	}
	if cur == nil {
		return nil, errors.New("base currency not found")
	}

	return graphdb.CurrencyFromDb(cur, c.Language()), nil
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
