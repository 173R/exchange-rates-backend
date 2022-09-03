package main

import (
	"github.com/wolframdeus/exchange-rates-backend/cmd"
	"time"
)

func main() {
	time.Local = time.UTC

	if err := cmd.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
