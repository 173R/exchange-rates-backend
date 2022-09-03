package http

import (
	"fmt"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/wolframdeus/exchange-rates-backend/configs"
	"github.com/wolframdeus/exchange-rates-backend/internal/db"
	http "github.com/wolframdeus/exchange-rates-backend/internal/http/middlewares"
	curreppkg "github.com/wolframdeus/exchange-rates-backend/internal/repositories/currencies"
	exratesreppkg "github.com/wolframdeus/exchange-rates-backend/internal/repositories/exrates"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories/refsessions"
	ureppkg "github.com/wolframdeus/exchange-rates-backend/internal/repositories/users"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/auth"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/currencies"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/exrates"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/users"
)

// Run запускает HTTP-сервер проекта.
func Run() error {
	// Создаем инстанс DB.
	gormDb, err := db.NewByConfig()
	if err != nil {
		return err
	}

	// Создаем все необходимые сервисы.
	uSrv := users.NewService(ureppkg.NewRepository(gormDb))
	exratesSrv := exrates.NewService(exratesreppkg.NewRepository(gormDb))
	curSrv := currencies.NewService(curreppkg.NewRepository(gormDb))
	authSrv := auth.NewService(uSrv, refsessions.NewRepository(gormDb))

	// Создаем корневой обработчик Gin.
	app := gin.New()

	// Добавляем обработчик CORS запросов.
	app.Use(cors.AllowAll())

	// Добавляем обработчик для процессинга паник.
	app.Use(http.NewCustomRecoveryMiddleware())

	// Инициализируем стандартный обработчик Sentry.
	app.Use(sentrygin.New(sentrygin.Options{
		// Мы перевыбрасываем ошибку для того, чтобы прокинуть её в другие
		// обработчики.
		Repanic: true,
	}))

	// Добавляем обработчик GraphQL запросов.
	app.POST("/gql", http.NewGraphQLMiddleware(curSrv, uSrv, exratesSrv, authSrv))

	if err := app.Run(fmt.Sprintf("0.0.0.0:%d", configs.App.Port)); err != nil {
		return err
	}
	return nil
}
