package graph

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/services/auth"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/currencies"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/exrates"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/users"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type services struct {
	// Сервис связанный с авторизацией.
	Auth *auth.Service
	// Сервис для работы с валютами.
	Currencies *currencies.Service
	// Сервис для работы с пользователями.
	Users *users.Service
	// Сервис для работы с курсами обменов валют.
	ExchangeRates *exrates.Service
}

type Resolver struct {
	Services *services
}

// NewResolver возвращает указатель на новый экземпляр Resolver.
func NewResolver(
	curSrv *currencies.Service,
	uSrv *users.Service,
	exRatesSrv *exrates.Service,
	authSrv *auth.Service,
) *Resolver {
	return &Resolver{
		Services: &services{
			Auth:          authSrv,
			Currencies:    curSrv,
			Users:         uSrv,
			ExchangeRates: exRatesSrv,
		},
	}
}
