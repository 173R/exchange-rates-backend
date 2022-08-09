package sentry

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/wolframdeus/exchange-rates-backend/configs"
)

// GetClient возвращает ссылку на новый инициализированный экземпляр
// клиента Sentry.
func GetClient(options sentry.ClientOptions) (*sentry.Client, error) {
	h := sentry.CurrentHub().Clone()
	client, err := sentry.NewClient(options)
	if err != nil {
		return nil, err
	}
	h.BindClient(client)

	return client, nil
}

// GetClientByConfig возвращает клиент Sentry по конфигу проекта.
func GetClientByConfig() (*sentry.Client, error) {
	return GetClient(sentry.ClientOptions{
		Dsn:         configs.Sentry.Dsn,
		Environment: configs.Sentry.Env,
	})
}

// Init инициализирует глобальный хаб Sentry.
func Init() error {
	fmt.Printf("%+v\n", configs.Sentry)
	return sentry.Init(sentry.ClientOptions{
		Dsn:         configs.Sentry.Dsn,
		Debug:       true,
		Environment: configs.Sentry.Env,
	})
}
