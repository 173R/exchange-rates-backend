package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "config",
		Short: "Выводит текущую конфигурацию проекта.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CONFIG")
		},
	})
}
