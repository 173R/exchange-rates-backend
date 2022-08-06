package configs

import (
	"fmt"
	"github.com/spf13/viper"
	"strconv"
)

// Создает ошибку по общему формату.
func createError(key string, t string) error {
	return fmt.Errorf(`значение по ключу "%s" не является %s`, key, t)
}

// Возвращает AppEnv на основе переданного ключа.
func getAppEnv(key string) AppEnv {
	value := AppEnv(getString(key))

	switch value {
	case AppEnvLocal, AppEnvProduction:
		return value
	default:
		panic(createError(key, "AppEnv"))
	}
}

// Возвращает uint на основе переданного ключа. В случае, если значение
// некорректно, будет выброшена ошибка.
func getRequiredUint(key string) uint {
	str := viper.GetString(key)

	if intVal, err := strconv.Atoi(str); err == nil {
		if intVal >= 0 {
			return uint(intVal)
		}
	}
	panic(createError(key, "uint"))
}

// Возвращает string на основе переданного ключа.
func getString(key string) string {
	return viper.GetString(key)
}

// Возвращает string на основе переданного ключа. В случае, если значение
// некорректно, будет выброшена ошибка.
func getRequiredString(key string) string {
	str := viper.GetString(key)

	if len(str) > 0 {
		return str
	}
	panic(createError(key, "string"))
}

// Возвращает значение по ключу предполагаю, что в нем указан порт.
func getPort(key string) uint {
	return getRequiredUint(key)
}
