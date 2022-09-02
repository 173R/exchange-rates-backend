package configs

type jwt struct {
	// Ключ, который используется для формирования JWT.
	Secret []byte
}

// Jwt содержит информацию касательно JSON Web Token-ов.
var Jwt *jwt

func init() {
	InitViper()
	Jwt = &jwt{
		Secret: []byte(getRequiredString("JWT_SECRET")),
	}
}
