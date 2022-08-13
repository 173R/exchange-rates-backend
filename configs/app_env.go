package configs

const (
	AppEnvLocal      AppEnv = "local"
	AppEnvStaging    AppEnv = "staging"
	AppEnvProduction AppEnv = "production"
)

// AppEnv описывает среду запуска приложения.
type AppEnv string
