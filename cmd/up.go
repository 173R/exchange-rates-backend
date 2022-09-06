package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/wolframdeus/exchange-rates-backend/configs"
	"github.com/wolframdeus/exchange-rates-backend/internal/db"
	"github.com/wolframdeus/exchange-rates-backend/internal/http"
	"github.com/wolframdeus/exchange-rates-backend/internal/sentry"
)

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "Запускает HTTP-сервер проекта вместе со всеми зависимостями.",
		Run: func(cmd *cobra.Command, args []string) {
			// Запускаем миграции.
			if err := db.RunMigrations(); err != nil {
				panic(err)
			}

			// Инициализируем Sentry.
			if err := sentry.InitByConfig(); err != nil {
				panic(err)
			}

			// Включаем релиз-режим, если режим отладки не включён.
			if !configs.App.Debug {
				gin.SetMode("release")
			}

			// Запускаем http-сервер.
			if err := http.Run(); err != nil {
				panic(err)
			}
		},
	})
}
