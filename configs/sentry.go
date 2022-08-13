package configs

type sentry struct {
	// Адрес для отправки событий.
	Dsn string
}

// Sentry содержит информацию касательно Sentry.
var Sentry *sentry

func init() {
	InitViper()
	Sentry = &sentry{
		Dsn: getString("SENTRY_DSN"),
	}
}
