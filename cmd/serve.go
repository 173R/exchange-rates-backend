package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wolframdeus/exchange-rates-backend/internal/http"
)

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "serve",
		Short: "Запускает HTTP-сервер проекта.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := http.Run(); err != nil {
				panic(err)
			}
		},
	})
}
