package http

import (
	"fmt"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/wolframdeus/exchange-rates-backend/configs"
	"github.com/wolframdeus/exchange-rates-backend/internal/db"
	http "github.com/wolframdeus/exchange-rates-backend/internal/http/middlewares"
)

type CurrencyModel struct {
	Id       string `json:"id"`
	TitleKey string `json:"title_key"`
	Sign     string `json:"sign"`
}

// Run запускает HTTP-сервер проекта.
func Run() error {
	// Создаем инстанс DB.
	gormDb, err := db.NewByConfig()
	if err != nil {
		return err
	}

	// Создаем корневой обработчик Gin.
	app := gin.New()

	// Добавляем обработчик CORS запросов.
	app.Use(cors.AllowAll())

	// Добавляем обработчик для процессинга паник.
	app.Use(http.NewCustomRecoveryMiddleware())

	// Внедряем в контекст запроса необходимые сервисы.
	app.Use(http.NewContextConfigMiddleware(gormDb))

	// Инициализируем стандартный обработчик Sentry.
	app.Use(sentrygin.New(sentrygin.Options{
		// Мы перевыбрасываем ошибку для того, чтобы прокинуть её в другие
		// обработчики.
		Repanic: true,
	}))

	// Добавляем обработчик GraphQL запросов.
	app.POST("/gql", http.NewGraphQLMiddleware())

	if err := app.Run(fmt.Sprintf("0.0.0.0:%d", configs.App.Port)); err != nil {
		return err
	}
	return nil
}
