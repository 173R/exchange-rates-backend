package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/wolframdeus/exchange-rates-backend/graph/generated"
	"github.com/wolframdeus/exchange-rates-backend/graph/model"
	ctxpkg "github.com/wolframdeus/exchange-rates-backend/internal/context"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
)

// Currencies is the resolver for the currencies field.
func (r *queryResolver) Currencies(ctx context.Context) ([]*model.Currency, error) {
	c := ctxpkg.NewGraph(&ctx)

	// Получаем список валют.
	currencies, err := c.Services.Currencies.FindAll()
	if err != nil {
		return nil, err
	}

	// Получаем язык для того, чтобы перевести наименования валют.
	lang := c.GetLanguage()

	// Получаем все ключи переводов.
	titleKeys := make([]models.TranslationId, len(currencies))

	for i, c := range currencies {
		titleKeys[i] = c.TitleTranslationId
	}

	// Получаем список всех переводов.
	var trlMap map[models.TranslationId]*models.Translation
	translations, err := c.Services.Translations.FindByIds(titleKeys)

	if err == nil {
		trlMap = make(map[models.TranslationId]*models.Translation, len(translations))

		for _, t := range translations {
			tValue := t
			trlMap[t.Id] = &tValue
		}
	} else {
		// TODO: Залогировать ошибку.
	}

	res := make([]*model.Currency, len(currencies))
	for i, c := range currencies {
		title := string(c.TitleTranslationId)

		// Переводы могли быть не найдены, поэтому необходимо совершить эту
		// проверку.
		if trlMap != nil {
			// Конкретный перевод также мог быть не найден.
			if t, ok := trlMap[c.TitleTranslationId]; ok {
				title = t.Translate(lang)
			}
		}

		res[i] = &model.Currency{
			ID:    string(c.Id),
			Sign:  c.Sign,
			Title: title,
		}
	}

	return res, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
