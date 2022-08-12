package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wolframdeus/exchange-rates-backend/internal/db"
)

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "seed",
		Short: "Запускает сиды проекта.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := db.RunSeeds(); err != nil {
				panic(err)
			}
		},
	})
}
