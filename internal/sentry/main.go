package sentry

import (
	"github.com/getsentry/sentry-go"
	"github.com/wolframdeus/exchange-rates-backend/configs"
)

// InitByConfig инициализирует глобальный хаб Sentry.
func InitByConfig() error {
	return sentry.Init(sentry.ClientOptions{
		Dsn:         configs.Sentry.Dsn,
		Debug:       configs.App.Debug,
		Environment: string(configs.App.Env),
	})
}
