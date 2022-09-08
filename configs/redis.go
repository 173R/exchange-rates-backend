package configs

type redis struct {
	// Хост БД.
	Host string
	// Порт БД.
	Port uint
	// Пароль БД.
	Pass string
}

// Redis содержит информацию касательно подключения к БД Redis.
var Redis *redis

func init() {
	InitViper()

	// Получаем хост.
	host := getString("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}

	// Получаем порт.
	port := getUint("REDIS_PORT")
	if port == 0 {
		port = 6379
	}

	Redis = &redis{
		Host: host,
		Port: port,
		Pass: getRequiredString("REDIS_PASS"),
	}
}
