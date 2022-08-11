package configs

import (
	"github.com/spf13/viper"
	"sync"
)

var once sync.Once

// FIXME: Мы установили значения по умолчанию в документации, но не в коде!

func InitViper() {
	once.Do(func() {
		// Указываем путь к конфиг файл.
		viper.SetConfigFile(".env")
		viper.AllowEmptyEnv(true)

		// Читаем конфиг из переменных окружения + .env файла.
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				return
			}
			panic(err)
		}
	})
}
