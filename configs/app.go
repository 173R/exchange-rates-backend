package configs

type app struct {
	// Текущая среда запуска приложения.
	Env AppEnv
	// Порт, на котором запускается HTTP-сервер.
	Port uint
}

// App содержит основную информацию о конфигурации проекта.
var App *app

func init() {
	InitViper()
	App = &app{
		Env:  getAppEnv("APP_ENV"),
		Port: getPort("PORT"),
	}
}
