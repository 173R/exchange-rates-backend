package configs

type app struct {
	// Включён ли режим отладки.
	Debug bool
	// Порт, на котором запускается HTTP-сервер.
	Port uint
}

// App содержит основную информацию о конфигурации проекта.
var App *app

func init() {
	InitViper()
	App = &app{
		Debug: getBoolean("DEBUG"),
		Port:  getPort("PORT"),
	}
}
