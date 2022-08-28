package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/wolframdeus/exchange-rates-backend/internal/db"
	exratesreppkg "github.com/wolframdeus/exchange-rates-backend/internal/repositories/exrates"
	exratessrvpkg "github.com/wolframdeus/exchange-rates-backend/internal/services/exrates"
)

func main() {
	d, _ := db.NewByConfig()
	rep := exratesreppkg.NewRepository(d)
	srv := exratessrvpkg.NewService(rep)
	spew.Dump(srv.FindPrevDayDiff("USD"))
}
