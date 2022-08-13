package launchparams

import "encoding/json"

type Params struct {
	// Локаль пользователя.
	Language string `json:"language"`
}

// Derive извлекает параметры запуска из какой-либо строки.
func Derive(value string) (*Params, error) {
	var j Params

	if err := json.Unmarshal([]byte(value), &j); err != nil {
		return nil, err
	}
	return &j, nil
}
