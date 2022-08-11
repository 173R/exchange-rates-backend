package http

import (
	"fmt"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/wolframdeus/exchange-rates-backend/configs"
	"github.com/wolframdeus/exchange-rates-backend/internal/context"
	http "github.com/wolframdeus/exchange-rates-backend/internal/http/middlewares"
	"github.com/wolframdeus/exchange-rates-backend/internal/sentry"
)

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

	// Создаем корневой обработчик Gin.
	app := gin.New()

	// Добавляем обработчик для процессинга паник.
	app.Use(http.NewCustomRecoveryMiddleware())

	// Инициализируем стандартный обработчик Sentry.
	app.Use(sentrygin.New(sentrygin.Options{
		// Мы перевыбрасываем ошибку для того, чтобы прокинуть её в другие
		// обработчики.
		Repanic: true,
	}))

	app.GET("/message", context.NewHandler(func(c *context.Context) {
		c.Send("its ok!")
	}))

	app.GET("/panic", context.NewHandler(func(c *context.Context) {
		panic("something went wrong with panic")
	}))

	if err := app.Run(fmt.Sprintf("0.0.0.0:%d", configs.App.Port)); err != nil {
		return err
	}
	return nil
}
