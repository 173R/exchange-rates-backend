package main

import "github.com/wolframdeus/exchange-rates-backend/cmd"

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
