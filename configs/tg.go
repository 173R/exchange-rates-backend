package configs

type tg struct {
	// Секретный ключ бота.
	SecretKey string
}

// Tg содержит информацию касательно Telegram.
var Tg *tg

func init() {
	InitViper()
	Tg = &tg{
		SecretKey: getRequiredString("TG_SECRET_KEY"),
	}
}
