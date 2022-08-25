package configs

type oer struct {
	// Ключ для работы с API Open Exchange Rates.
	AppId string
}

// OER содержит информацию касающуюся Open Exchange Rates.
var OER *oer

func init() {
	InitViper()
	OER = &oer{
		AppId: getRequiredString("OER_APP_ID"),
	}
}
