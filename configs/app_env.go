package configs

const (
	AppEnvLocal      AppEnv = "local"
	AppEnvProduction AppEnv = "production"
)

// AppEnv описывает среду запуска приложения.
type AppEnv string
