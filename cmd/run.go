package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wolframdeus/exchange-rates-backend/internal/db"
	"github.com/wolframdeus/exchange-rates-backend/internal/http"
)

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "Запускает HTTP-сервер проекта вместе с миграциями.",
		Run: func(cmd *cobra.Command, args []string) {
			// Запускаем миграции.
			if err := db.RunMigrations(); err != nil {
				panic(err)
			}

			// Запускаем http-сервер.
			if err := http.Run(); err != nil {
				panic(err)
			}
		},
	})
}
