package configs

type db struct {
	// Хост, на котором расположена БД.
	Host string
	// Порт, на котором расположена БД.
	Port uint
	// Наименование БД.
	Name string
	// Пользователь БД.
	User string
	// Пароль пользователя БД.
	Pass string
}

// Db содержит информацию касательно базы данных.
var Db *db

func init() {
	InitViper()
	Db = &db{
		Host: getRequiredString("DB_HOST"),
		Name: getRequiredString("DB_NAME"),
		User: getRequiredString("DB_USER"),
		Pass: getRequiredString("DB_PASS"),
		Port: getPort("DB_PORT"),
	}
}
