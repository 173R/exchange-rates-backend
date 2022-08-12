package http

import (
	"fmt"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/wolframdeus/exchange-rates-backend/configs"
	"github.com/wolframdeus/exchange-rates-backend/internal/context"
	"github.com/wolframdeus/exchange-rates-backend/internal/db"
	http "github.com/wolframdeus/exchange-rates-backend/internal/http/middlewares"
	"github.com/wolframdeus/exchange-rates-backend/internal/sentry"
)

type CurrencyModel struct {
	Id       string `json:"id"`
	TitleKey string `json:"title_key"`
	Sign     string `json:"sign"`
}

// Run запускает HTTP-сервер проекта.
func Run() error {
	// FIXME: Использовать локальный клиент.
	//// Создаём клиент Sentry.
	//sentryClient, err := sentry.GetClientByConfig()
	//if err != nil {
	//	return err
	//}
	if err := sentry.Init(); err != nil {
		return err
	}

	// Включаем релиз-режим, если режим отладки не включён.
	if !configs.App.Debug {
		gin.SetMode("release")
	}

	// Создаем инстанс DB.
	gormDb, err := db.NewByConfig()
	if err != nil {
		return err
	}

	// Создаем корневой обработчик Gin.
	app := gin.New()

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

	app.GET("/currencies", context.NewHandler(func(c *context.Context) {
		// Получаем список всех валют.
		cur, err := c.Services.Currencies.FindAll()
		if err != nil {
			c.SendError(err)
			return
		}

		// Конвертируем полученные результаты к моделям.
		res := make([]CurrencyModel, len(cur))
		for i, c := range cur {
			res[i] = CurrencyModel{
				Id:       string(c.Id),
				TitleKey: c.TitleKey,
				Sign:     c.Sign,
			}
		}

		c.SendData(res)
	}))

	if err := app.Run(fmt.Sprintf("0.0.0.0:%d", configs.App.Port)); err != nil {
		return err
	}
	return nil
}
