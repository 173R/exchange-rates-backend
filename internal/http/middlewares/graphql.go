package http

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/gqlerror"
	ctxpkg "github.com/wolframdeus/exchange-rates-backend/internal/context"
	"github.com/wolframdeus/exchange-rates-backend/internal/customerrors"
	"github.com/wolframdeus/exchange-rates-backend/internal/graph"
	"github.com/wolframdeus/exchange-rates-backend/internal/graph/generated"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/auth"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/currencies"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/exrates"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/users"
	"strconv"
	"strings"
)

// NewGraphQLMiddleware возвращает новый обработчик для GraphQL запросов.
func NewGraphQLMiddleware(
	curSrv *currencies.Service,
	uSrv *users.Service,
	exRatesSrv *exrates.Service,
	authSrv *auth.Service,
) gin.HandlerFunc {
	h := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{
			Resolvers: graph.NewResolver(curSrv, uSrv, exRatesSrv, authSrv),
		}),
	)

	// Добавляем кастомный обработчик ошибок, который отправит их в Sentry.
	h.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		w := ctxpkg.MustNewGraph(ctx)

		// Применяем стандартный форматтер ошибки.
		errFormatted := graphql.DefaultErrorPresenter(ctx, err)

		// Логируем ошибку при помощи Sentry.
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelError)

			// Устанавливаем данные пользователя.
			if uid, err := w.UserId(); err == nil {
				scope.SetUser(sentry.User{
					ID: strconv.FormatInt(int64(uid), 10),
					// TODO: Указать адрес.
					IPAddress: "",
				})
			}

			// Пытаемся привести ошибку к нашему кастомному типу.
			var myErr *customerrors.Error
			var errText = err.Error()

			if errors.As(err, &myErr) {
				// Если ошибка действительно имеет наш кастомный тип, мы должны
				// добавить все хлебные крошки.
				breadcrumbsToScope(myErr, scope)

				// Подменяем сообщение ошибки на публичное.
				errFormatted.Message = myErr.PublicError()

				// Переопределяем текст ошибки.
				errText = myErr.Error()
			}
			sentry.CaptureException(errors.New(errText))
		})

		return errFormatted
	})

	return func(c *gin.Context) {
		// Извлекаем контекст запроса.
		ctx := c.Request.Context()

		// Помещаем в контекст кеш запроса.
		ctx = ctxpkg.ContextWithCache(ctx, uSrv)

		// Помещаем в контекст токен авторизации.
		ctx = context.WithValue(ctx, ctxpkg.KeyAuthToken, deriveAuthToken(c))

		// Переопределяем контекст запроса.
		c.Request = c.Request.WithContext(ctx)

		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Извлекает токен авторизации из запроса.
func deriveAuthToken(c *gin.Context) string {
	h := c.GetHeader("authorization")
	parts := strings.Split(h, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}

// Пытается представить переданную ошибку как нашу кастомную и в этом
// случае добавить в нужном порядке все хлебные крошки.
func breadcrumbsToScope(e error, scope *sentry.Scope) {
	if err, ok := e.(*customerrors.Error); ok {
		breadcrumbsToScope(err.Original, scope)

		if len(err.Breadcrumbs) > 0 {
			for _, b := range err.Breadcrumbs {
				scope.AddBreadcrumb(b, 1000)
			}
		}
	}
}
