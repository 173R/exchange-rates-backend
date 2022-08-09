package configs

type sentry struct {
	// Адрес для отправки событий.
	Dsn string
	// Среда приложения.
	Env string
}

// Sentry содержит информацию касательно Sentry.
var Sentry *sentry

func init() {
	InitViper()
	Sentry = &sentry{
		Dsn: getString("SENTRY_DSN"),
		Env: getString("SENTRY_ENV"),
	}
}
